import { useEffect, useRef, useState } from "react";
import { getJSON, sendJSON } from "../api";
import type { ChatProvider } from "../types";

type Msg = { role: "user" | "assistant"; content: string; off?: boolean };

export default function Chatbot() {
  const [open, setOpen] = useState(false);
  const [providers, setProviders] = useState<ChatProvider[]>([]);
  const [active, setActive] = useState("gemini");
  const [input, setInput] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [msgs, setMsgs] = useState<Msg[]>([
    {
      role: "assistant",
      content:
        "Ask about Yo Shandy Hermawan's profile, stack, projects, or engineering approach. I only answer from the portfolio data.",
    },
  ]);
  const scroller = useRef<HTMLDivElement>(null);

  useEffect(() => {
    getJSON<{ active: string; providers: ChatProvider[] }>("/api/chat/providers")
      .then((d) => {
        setProviders(d.providers);
        const firstOn = d.providers.find((p) => p.enabled)?.id;
        setActive(firstOn || d.active || "gemini");
      })
      .catch(() => undefined);
  }, []);

  useEffect(() => {
    scroller.current?.scrollTo({ top: scroller.current.scrollHeight });
  }, [msgs, open]);

  async function send() {
    const text = input.trim();
    if (!text || busy) return;
    setInput("");
    setError("");
    const next = [...msgs, { role: "user" as const, content: text }];
    setMsgs(next);
    setBusy(true);
    try {
      const history = next
        .filter((m) => m.role === "user" || m.content)
        .slice(-8)
        .map((m) => ({ role: m.role, content: m.content }));
      history.pop();
      const res = await sendJSON<{ answer: string; off_topic: boolean }>("/api/chat", "POST", {
        message: text,
        provider: active,
        history,
      });
      setMsgs((cur) => [...cur, { role: "assistant", content: res.answer, off: res.off_topic }]);
    } catch (e) {
      setError(e instanceof Error ? e.message : "chat failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <>
      <button className="chat-toggle" onClick={() => setOpen((v) => !v)}>
        {open ? "Close chat" : "Ask recruiter bot"}
      </button>
      {open && (
        <aside className="chat">
          <header>
            <div>
              <h2>Portfolio bot</h2>
              <div className="hint" style={{ padding: 0 }}>
                Grounded in porto.md · off-topic questions are filtered
              </div>
            </div>
            <select value={active} onChange={(e) => setActive(e.target.value)} aria-label="LLM">
              {providers.map((p) => (
                <option key={p.id} value={p.id} disabled={!p.enabled}>
                  {p.label}
                  {!p.enabled ? " · no key" : ""}
                </option>
              ))}
            </select>
          </header>
          <div className="msgs" ref={scroller}>
            {msgs.map((m, i) => (
              <div key={i} className={`msg ${m.role === "user" ? "user" : "bot"} ${m.off ? "off" : ""}`}>
                {m.content}
              </div>
            ))}
            {busy && <div className="msg bot">Thinking…</div>}
            {error && <div className="flash err">{error}</div>}
          </div>
          <div className="composer">
            <textarea
              value={input}
              placeholder="e.g. Why was the Omni-chan webhook separated?"
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  void send();
                }
              }}
            />
            <button className="primary" disabled={busy} onClick={() => void send()}>
              Send
            </button>
          </div>
        </aside>
      )}
    </>
  );
}
