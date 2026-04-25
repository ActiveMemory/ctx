//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"encoding/json"
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
// Returns ("", nil) when:
//   - the key file is missing (key was never generated),
//   - the encrypted file is missing (webhook never configured).
//
// Any resolver or I/O failure is propagated (including
// [errCtx.ErrDirNotDeclared]) so callers can distinguish
// "no context dir" from "no webhook configured" rather than
// being forced to treat them identically. [Send] treats any error
// as "no webhook, silently skip"; interactive callers (e.g.
// `ctx notify test`) can use [errors.Is] to surface a clearer
// message when the project is not set up yet.
//
// Returns:
//   - string: the decrypted webhook URL, or "" if not configured
//   - error: non-nil on any resolver failure or decryption failure;
//     missing key / encrypted file are silent
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

	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		if os.IsNotExist(loadErr) {
			return "", nil
		}
		return "", nil
	}

	ciphertext, readErr := io.SafeReadUserFile(encPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return "", nil
		}
		return "", nil
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

// Send fires a webhook notification. It is a silent noop when:
//   - no webhook URL is configured
//   - the event is not in the allowed list
//   - the HTTP request fails (fire-and-forget)
//
// Parameters:
//   - event: notification category (e.g. "relay", "nudge")
//   - message: short human-readable summary
//   - sessionID: Claude Code session ID (may be empty)
//   - detail: structured template reference (nil omits the field)
//
// Returns:
//   - error: Delivery error, or nil if sent successfully or silently skipped
func Send(event, message, sessionID string, detail *entity.TemplateRef) error {
	if !EventAllowed(event, rc.NotifyEvents()) {
		return nil
	}

	url, webhookErr := LoadWebhook()
	if webhookErr != nil || url == "" {
		return nil
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
		return nil
	}

	resp, postErr := PostJSON(url, body)
	if postErr != nil {
		return nil // fire-and-forget
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
