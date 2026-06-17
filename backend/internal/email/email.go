package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Sender is the interface every notification path in the platform uses.
// The real Resend-backed implementation and the local dev-outbox fallback
// both satisfy it, so callers never branch on mock mode themselves.
type Sender interface {
	Send(ctx context.Context, to, subject, htmlBody string) error
}

func NewSender(apiKey, fromAddr string, mock bool) Sender {
	if mock {
		return &devOutboxSender{dir: ".devmail"}
	}
	return &resendSender{apiKey: apiKey, fromAddr: fromAddr, client: &http.Client{Timeout: 10 * time.Second}}
}

type resendSender struct {
	apiKey   string
	fromAddr string
	client   *http.Client
}

func (s *resendSender) Send(ctx context.Context, to, subject, htmlBody string) error {
	payload := map[string]any{
		"from":    fmt.Sprintf("Team1 Blog <%s>", s.fromAddr),
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// devOutboxSender writes the rendered email to disk so the full notification
// flow can be exercised locally without a Resend API key.
type devOutboxSender struct {
	dir string
}

func (s *devOutboxSender) Send(ctx context.Context, to, subject, htmlBody string) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s_%s_%s.html", time.Now().Format("20060102T150405"), sanitize(to), sanitize(subject))
	path := filepath.Join(s.dir, filename)
	fmt.Printf("[devmail] to=%s subject=%q -> %s\n", to, subject, path)
	return os.WriteFile(path, []byte(htmlBody), 0o644)
}

func sanitize(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			out = append(out, r)
		default:
			out = append(out, '-')
		}
	}
	if len(out) > 60 {
		out = out[:60]
	}
	return string(out)
}
