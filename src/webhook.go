package main

// Webhook fired after every successful push (StatusPushed) so the
// IKAROS bridge (api/idig-sync) can pull just the affected trench
// without polling. Configured via two env vars; if either is unset
// the webhook is disabled and SyncTrench behaves exactly as before.
//
//	IDIG_WEBHOOK_URL    e.g. http://api:3001/idig-sync/webhook
//	IDIG_WEBHOOK_SECRET hex/random shared secret (HMAC-SHA256 key)
//
// Delivery is fire-and-forget on a goroutine: failures are logged
// but never block the Sync response. The bridge has a periodic
// full-sync fallback, so a dropped webhook just delays one trench.

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type WebhookConfig struct {
	URL    string
	Secret []byte
	client *http.Client
}

// LoadWebhookConfig reads env once. Returns nil if the webhook is
// not configured (either var missing); callers should treat nil as
// "do nothing".
func LoadWebhookConfig() *WebhookConfig {
	url := os.Getenv("IDIG_WEBHOOK_URL")
	secret := os.Getenv("IDIG_WEBHOOK_SECRET")
	if url == "" || secret == "" {
		return nil
	}
	return &WebhookConfig{
		URL:    url,
		Secret: []byte(secret),
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

type pushedPayload struct {
	Project string `json:"project"`
	Trench  string `json:"trench"`
	Commit  string `json:"commit"`
	User    string `json:"user"`
	Ts      int64  `json:"ts"`
}

// FirePushed sends the webhook on a goroutine. Safe to call when the
// receiver is nil.
func (w *WebhookConfig) FirePushed(project, trench, commit, user string) {
	if w == nil {
		return
	}
	payload := pushedPayload{
		Project: project,
		Trench:  trench,
		Commit:  commit,
		User:    user,
		Ts:      time.Now().UnixMilli(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("webhook: marshal: %v", err)
		return
	}

	mac := hmac.New(sha256.New, w.Secret)
	mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	go func() {
		req, err := http.NewRequest(http.MethodPost, w.URL, bytes.NewReader(body))
		if err != nil {
			log.Printf("webhook: build request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-IDig-Signature", sig)

		resp, err := w.client.Do(req)
		if err != nil {
			log.Printf("webhook: POST %s failed: %v", w.URL, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			log.Printf("webhook: POST %s -> %d", w.URL, resp.StatusCode)
		}
	}()
}
