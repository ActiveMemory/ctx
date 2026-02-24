//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

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

// testCmd returns the "ctx notify test" subcommand.
func testCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Send a test notification",
		Long:  `Sends a test notification to the configured webhook and reports the HTTP status.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTest(cmd)
		},
	}
}

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

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	// Check event filter — but for test we bypass and send directly
	if !notify.EventAllowed("test", rc.NotifyEvents()) {
		cmd.Println("Note: event \"test\" is filtered by your .ctxrc notify.events config.")
		cmd.Println("Sending anyway for testing purposes.")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(body)) //nolint:gosec // URL is user-configured via encrypted storage
	if err != nil {
		return fmt.Errorf("send test notification: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	cmd.Println(fmt.Sprintf("Webhook responded: HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode)))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		cmd.Println("Webhook is working " + config.FileNotifyEnc)
	}

	return nil
}
