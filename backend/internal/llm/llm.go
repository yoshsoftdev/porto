package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Provider string

const (
	Gemini Provider = "gemini"
	GPT    Provider = "gpt"
	Grok   Provider = "grok"
	Claude Provider = "claude"
)

func All() []Provider {
	return []Provider{Gemini, GPT, Grok, Claude}
}

func Valid(name string) bool {
	switch Provider(name) {
	case Gemini, GPT, Grok, Claude:
		return true
	}
	return false
}

func DefaultModel(p Provider) string {
	switch p {
	case Gemini:
		return "gemini-2.0-flash"
	case GPT:
		return "gpt-4o-mini"
	case Grok:
		return "grok-3-mini"
	case Claude:
		return "claude-3-5-haiku-latest"
	default:
		return ""
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Complete(ctx context.Context, provider Provider, model, apiKey, system string, history []Message) (string, error) {
	if strings.TrimSpace(apiKey) == "" {
		return "", fmt.Errorf("api key for %s is not configured", provider)
	}
	if model == "" {
		model = DefaultModel(provider)
	}
	switch provider {
	case Gemini:
		return gemini(ctx, model, apiKey, system, history)
	case GPT:
		return openAICompat(ctx, "https://api.openai.com/v1/chat/completions", apiKey, model, system, history)
	case Grok:
		return openAICompat(ctx, "https://api.x.ai/v1/chat/completions", apiKey, model, system, history)
	case Claude:
		return claude(ctx, model, apiKey, system, history)
	default:
		return "", fmt.Errorf("unknown provider %s", provider)
	}
}

func gemini(ctx context.Context, model, key, system string, history []Message) (string, error) {
	type part struct {
		Text string `json:"text"`
	}
	type content struct {
		Role  string `json:"role,omitempty"`
		Parts []part `json:"parts"`
	}
	var contents []content
	for _, m := range history {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, content{Role: role, Parts: []part{{Text: m.Content}}})
	}
	body := map[string]any{
		"system_instruction": content{Parts: []part{{Text: system}}},
		"contents":           contents,
		"generationConfig": map[string]any{
			"temperature":     0.3,
			"maxOutputTokens": 1024,
		},
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, key)
	raw, err := postJSON(ctx, url, nil, body)
	if err != nil {
		return "", err
	}
	var resp struct {
		Candidates []struct {
			Content struct {
				Parts []part `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("gemini: %s", resp.Error.Message)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned an empty response")
	}
	var b strings.Builder
	for _, p := range resp.Candidates[0].Content.Parts {
		b.WriteString(p.Text)
	}
	return strings.TrimSpace(b.String()), nil
}

func openAICompat(ctx context.Context, url, key, model, system string, history []Message) (string, error) {
	msgs := []map[string]string{{"role": "system", "content": system}}
	for _, m := range history {
		role := m.Role
		if role != "assistant" {
			role = "user"
		}
		msgs = append(msgs, map[string]string{"role": role, "content": m.Content})
	}
	body := map[string]any{
		"model":       model,
		"messages":    msgs,
		"temperature": 0.3,
		"max_tokens":  1024,
	}
	raw, err := postJSON(ctx, url, map[string]string{
		"Authorization": "Bearer " + key,
	}, body)
	if err != nil {
		return "", err
	}
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("%s", resp.Error.Message)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty model response")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func claude(ctx context.Context, model, key, system string, history []Message) (string, error) {
	var msgs []map[string]string
	for _, m := range history {
		role := m.Role
		if role != "assistant" {
			role = "user"
		}
		msgs = append(msgs, map[string]string{"role": role, "content": m.Content})
	}
	body := map[string]any{
		"model":      model,
		"max_tokens": 1024,
		"system":     system,
		"messages":   msgs,
	}
	raw, err := postJSON(ctx, "https://api.anthropic.com/v1/messages", map[string]string{
		"x-api-key":         key,
		"anthropic-version": "2023-06-01",
	}, body)
	if err != nil {
		return "", err
	}
	var resp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("claude: %s", resp.Error.Message)
	}
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("claude returned an empty response")
	}
	return strings.TrimSpace(resp.Content[0].Text), nil
}

func postJSON(ctx context.Context, url string, headers map[string]string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 45 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		msg := strings.TrimSpace(string(raw))
		if len(msg) > 400 {
			msg = msg[:400]
		}
		return nil, fmt.Errorf("llm http %d: %s", res.StatusCode, msg)
	}
	return raw, nil
}
