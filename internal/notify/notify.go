//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cfgCrypto "github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHTTP "github.com/ActiveMemory/ctx/internal/config/http"
	"github.com/ActiveMemory/ctx/internal/config/project"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// LoadWebhook reads and decrypts the webhook URL from .context/.notify.enc.
//
// The webhook is "configured" exactly when .context/.notify.enc exists.
// Its absence is the one silent "not configured" signal: LoadWebhook
// returns ("", nil). Once the encrypted file exists, anything that then
// prevents decryption is a real problem and is propagated: a missing,
// unreadable, or invalid key (e.g. a project-local key absent in a git
// worktree), an unreadable .notify.enc, a decryption failure (wrong
// key), or a resolver error such as [errCtx.ErrNoCtxHere]. This lets
// callers distinguish "not configured" (silent) from "configured but
// broken" (surface it). [Send] warns on the latter; interactive callers
// (e.g. `ctx hook notify test`) report it directly.
//
// Returns:
//   - string: the decrypted webhook URL, or "" if not configured
//   - error: non-nil when a configured webhook cannot be decrypted
func LoadWebhook() (string, error) {
	kp, kpErr := rc.KeyPath()
	if kpErr != nil {
		return "", kpErr
	}
	ctxDir, pathErr := rc.ContextDir()
	if pathErr != nil {
		return "", pathErr
	}
	encPath := filepath.Join(ctxDir, cfgCrypto.NotifyEnc)

	// A missing .notify.enc is the only silent "not configured" case.
	// os.Stat returns an unwrapped *fs.PathError, so this not-exist
	// check is reliable regardless of how downstream library errors are
	// wrapped (crypto.LoadKey wraps through the text registry, on which
	// neither os.IsNotExist nor errors.Is is dependable). Once the
	// encrypted file exists the webhook IS configured, so a missing,
	// unreadable, or invalid key — and any decrypt failure — is surfaced
	// rather than mistaken for "no webhook".
	if _, statErr := os.Stat(encPath); statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return "", nil // webhook never configured
		}
		return "", statErr
	}

	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return "", loadErr // configured, but key missing/unreadable/invalid
	}
	ciphertext, readErr := io.SafeReadUserFile(encPath)
	if readErr != nil {
		return "", readErr // enc present but unreadable
	}
	plaintext, decryptErr := crypto.Decrypt(key, ciphertext)
	if decryptErr != nil {
		return "", decryptErr
	}

	return string(plaintext), nil
}

// SaveWebhook encrypts and writes the webhook URL to .context/.notify.enc.
//
// If the scratchpad key does not exist, it is generated and saved first.
//
// Parameters:
//   - url: the webhook endpoint to store
//
// Returns:
//   - error: non-nil if key generation, encryption, or file write fails
func SaveWebhook(url string) error {
	kp, kpErr := rc.KeyPath()
	if kpErr != nil {
		return kpErr
	}
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return ctxErr
	}
	encPath := filepath.Join(ctxDir, cfgCrypto.NotifyEnc)

	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		// Key doesn't exist: generate one.
		var genErr error
		key, genErr = crypto.GenerateKey()
		if genErr != nil {
			return genErr
		}
		if mkdirErr := io.SafeMkdirAll(
			filepath.Dir(kp), fs.PermKeyDir,
		); mkdirErr != nil {
			return mkdirErr
		}
		if saveErr := crypto.SaveKey(kp, key); saveErr != nil {
			return saveErr
		}
	}

	ciphertext, encryptErr := crypto.Encrypt(key, []byte(url))
	if encryptErr != nil {
		return encryptErr
	}

	return io.SafeWriteFile(encPath, ciphertext, fs.PermSecret)
}

// EventAllowed reports whether the given event passes the filter.
//
// A nil or empty allowed list means no events pass (opt-in only).
//
// Parameters:
//   - event: the event name to check
//   - allowed: list of permitted event names
//
// Returns:
//   - bool: true if event appears in the allowed list
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

// Send fires a webhook notification. It is a silent noop only when
// delivery is not expected:
//   - the event is not in the allowed list (not subscribed), or
//   - no webhook URL is configured.
//
// When a webhook IS configured but cannot be delivered — an
// unreadable or wrong key, a decrypt failure (e.g. a project-local
// key absent in a git worktree), a marshal error, or an HTTP failure
// — Send emits a non-fatal warning to stderr and returns nil. It
// never returns a delivery error (fire-and-forget), but it is never
// silent about a real failure: a webhook the user set up that drops
// without a trace reads as "working" when it is not.
//
// Parameters:
//   - event: notification category (e.g. "relay", "nudge")
//   - message: short human-readable summary
//   - sessionID: Claude Code session ID (may be empty)
//   - detail: structured template reference (nil omits the field)
//
// Returns:
//   - error: always nil; failures are warned, not returned
func Send(event, message, sessionID string, detail *entity.TemplateRef) error {
	if !EventAllowed(event, rc.NotifyEvents()) {
		return nil
	}

	url, webhookErr := LoadWebhook()
	if webhookErr != nil {
		// Configured but undeliverable (wrong/absent key in a
		// worktree, unreadable key, or decrypt failure). Surface it,
		// but stay non-fatal (fire-and-forget).
		logWarn.Warn(cfgWarn.NotifyWebhookLoad, webhookErr)
		return nil
	}
	if url == "" {
		return nil // not configured: legitimate silent no-op
	}

	projectName := project.FallbackName
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		projectName = filepath.Base(cwd)
	} else {
		logWarn.Warn(cfgWarn.Getwd, cwdErr)
	}

	payload := entity.NewNotifyPayload(
		event, message, sessionID, projectName, detail,
	)

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		logWarn.Warn(cfgWarn.NotifyWebhookMarshal, marshalErr)
		return nil
	}

	resp, postErr := PostJSON(url, body)
	if postErr != nil {
		// Delivery failed: fire-and-forget, but no longer silent.
		logWarn.Warn(cfgWarn.NotifyWebhookPost, postErr)
		return nil
	}
	if closeErr := resp.Body.Close(); closeErr != nil {
		logWarn.Warn(cfgWarn.CloseResponse, closeErr)
	}

	return nil
}

// PostJSON sends a JSON payload to a webhook URL and returns the response.
// The URL is always user-configured via encrypted storage.
//
// Parameters:
//   - url: webhook endpoint.
//   - body: JSON-encoded payload bytes.
//
// Returns:
//   - *http.Response: the HTTP response (caller must close Body).
//   - error: on HTTP failure.
func PostJSON(url string, body []byte) (*http.Response, error) {
	return io.SafePost(
		url, cfgHTTP.MimeJSON, body,
		cfgHTTP.WebhookTimeout*time.Second)
}

// MaskURL shows the scheme + host and masks everything after the path start.
//
// Parameters:
//   - url: full webhook URL.
//
// Returns:
//   - string: masked URL safe for display.
func MaskURL(url string) string {
	count := 0
	for i, c := range url {
		if c == cfgHTTP.PathSep {
			count++
			if count == cfgHTTP.MaskAfterSlash {
				return url[:i] + cfgHTTP.PathSepStr + cfgHTTP.MaskSuffix
			}
		}
	}
	if len(url) > cfgHTTP.MaskMaxLen {
		return url[:cfgHTTP.MaskMaxLen] + cfgHTTP.MaskSuffix
	}
	return url
}
