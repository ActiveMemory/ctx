//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify provides fire-and-forget webhook notifications.
//
// The webhook URL is stored encrypted in .context/.notify.enc using the
// same AES-256-GCM key as the scratchpad (.context/.context.key).
// When no webhook is configured, all operations are silent noops.
package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Event     string `json:"event"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	Timestamp string `json:"timestamp"`
	Project   string `json:"project"`
}

// LoadWebhook reads and decrypts the webhook URL from .context/.notify.enc.
//
// Returns ("", nil) if either the key file or encrypted file is missing
// (silent noop — webhook not configured).
func LoadWebhook() (string, error) {
	contextDir := rc.ContextDir()
	config.MigrateKeyFile(contextDir)
	keyPath := filepath.Join(contextDir, config.FileContextKey)
	encPath := filepath.Join(contextDir, config.FileNotifyEnc)

	key, err := crypto.LoadKey(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", nil
	}

	ciphertext, err := os.ReadFile(encPath) //nolint:gosec // project-local path
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", nil
	}

	plaintext, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SaveWebhook encrypts and writes the webhook URL to .context/.notify.enc.
//
// If the scratchpad key does not exist, it is generated and saved first.
func SaveWebhook(url string) error {
	contextDir := rc.ContextDir()
	config.MigrateKeyFile(contextDir)
	keyPath := filepath.Join(contextDir, config.FileContextKey)
	encPath := filepath.Join(contextDir, config.FileNotifyEnc)

	key, err := crypto.LoadKey(keyPath)
	if err != nil {
		// Key doesn't exist — generate one.
		key, err = crypto.GenerateKey()
		if err != nil {
			return err
		}
		if saveErr := crypto.SaveKey(keyPath, key); saveErr != nil {
			return saveErr
		}
	}

	ciphertext, err := crypto.Encrypt(key, []byte(url))
	if err != nil {
		return err
	}

	return os.WriteFile(encPath, ciphertext, config.PermSecret)
}

// EventAllowed reports whether the given event passes the filter.
//
// A nil or empty allowed list means no events pass (opt-in only).
func EventAllowed(event string, allowed []string) bool {
	if len(allowed) == 0 {
		return false
	}
	for _, e := range allowed {
		if e == event {
			return true
		}
	}
	return false
}

// Send fires a webhook notification. It is a silent noop when:
//   - no webhook URL is configured
//   - the event is not in the allowed list
//   - the HTTP request fails (fire-and-forget)
func Send(event, message, sessionID string) error {
	if !EventAllowed(event, rc.NotifyEvents()) {
		return nil
	}

	url, err := LoadWebhook()
	if err != nil || url == "" {
		return nil
	}

	project := "unknown"
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		project = filepath.Base(cwd)
	}

	payload := Payload{
		Event:     event,
		Message:   message,
		SessionID: sessionID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   project,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(body)) //nolint:gosec // URL is user-configured via encrypted storage
	if err != nil {
		return nil // fire-and-forget
	}
	_ = resp.Body.Close()

	return nil
}
