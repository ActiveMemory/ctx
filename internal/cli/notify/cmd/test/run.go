//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runTest sends a test notification to the configured webhook.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on webhook load or HTTP failure
func runTest(cmd *cobra.Command) error {
	url, err := notify.LoadWebhook()
	if err != nil {
		return fmt.Errorf("load webhook: %w", err)
	}
	if url == "" {
		cmd.Println("No webhook configured. Run: ctx notify setup")
		return nil
	}

	project := "unknown"
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		project = filepath.Base(cwd)
	}

	payload := notify.Payload{
		Event:     "test",
		Message:   "Test notification from ctx",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   project,
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("marshal payload: %w", marshalErr)
	}

	// Check event filter — but for test we bypass and send directly
	if !notify.EventAllowed("test", rc.NotifyEvents()) {
		cmd.Println("Note: event \"test\" is filtered by your .ctxrc notify.events config.")
		cmd.Println("Sending anyway for testing purposes.")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, postErr := client.Post(url, "application/json", bytes.NewReader(body)) //nolint:gosec // URL is user-configured via encrypted storage
	if postErr != nil {
		return fmt.Errorf("send test notification: %w", postErr)
	}
	defer func() { _ = resp.Body.Close() }()

	cmd.Println(fmt.Sprintf("Webhook responded: HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode)))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		cmd.Println("Webhook is working " + config.FileNotifyEnc)
	}

	return nil
}
