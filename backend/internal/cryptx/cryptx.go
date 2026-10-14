package cryptx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func Seal(key, plaintext string) (string, error) {
	if strings.TrimSpace(plaintext) == "" {
		return "", nil
	}
	if strings.TrimSpace(key) == "" {
		return "plain:" + plaintext, nil
	}
	block, err := aes.NewCipher(derive(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return "enc:" + base64.StdEncoding.EncodeToString(out), nil
}

func Open(key, stored string) (string, error) {
	if stored == "" {
		return "", nil
	}
	if v, ok := strings.CutPrefix(stored, "plain:"); ok {
		return v, nil
	}
	raw, ok := strings.CutPrefix(stored, "enc:")
	if !ok {
		return stored, nil
	}
	if strings.TrimSpace(key) == "" {
		return "", errors.New("ENCRYPTION_KEY is required to read stored secrets")
	}
	bin, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(derive(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(bin) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ct := bin[:gcm.NonceSize()], bin[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

func Masked(stored string) bool {
	return strings.TrimSpace(stored) != ""
}

func derive(key string) []byte {
	sum := sha256.Sum256([]byte(key))
	return sum[:]
}
