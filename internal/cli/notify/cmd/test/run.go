//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/crypto"
	errNotify "github.com/ActiveMemory/ctx/internal/err/notify"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeNotify "github.com/ActiveMemory/ctx/internal/write/notify"
)

// runTest sends a test notification to the configured webhook.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on webhook load or HTTP failure
func runTest(cmd *cobra.Command) error {
	url, loadErr := notify.LoadWebhook()
	if loadErr != nil {
		return errNotify.LoadWebhook(loadErr)
	}
	if url == "" {
		writeNotify.TestNoWebhook(cmd)
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
		return errNotify.MarshalPayload(marshalErr)
	}

	if !notify.EventAllowed("test", rc.NotifyEvents()) {
		writeNotify.TestFiltered(cmd)
	}

	resp, postErr := notify.PostJSON(url, body)
	if postErr != nil {
		return errNotify.SendNotification(postErr)
	}
	defer func() { _ = resp.Body.Close() }()

	writeNotify.TestResult(cmd, resp.StatusCode, crypto.NotifyEnc)

	return nil
}
