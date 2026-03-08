//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"net/http"

	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/ActiveMemory/ctx/internal/write/io"
	"github.com/spf13/cobra"
)

// SetupPrompt prints the interactive webhook URL prompt.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func SetupPrompt(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Print(config.tplSetupPrompt)
}

// SetupDone prints the success block after saving a webhook:
// configured URL (masked) and encrypted file path.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - maskedURL: masked webhook URL for display.
//   - encPath: encrypted file path.
func SetupDone(cmd *cobra.Command, maskedURL, encPath string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplSetupDone, maskedURL, encPath)
}

// TestNoWebhook prints the message when no webhook is configured.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func TestNoWebhook(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplTestNoWebhook)
}

// TestFiltered prints the notice when the test event is filtered.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func TestFiltered(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplTestFiltered)
}

// TestResult prints the webhook test response block: status line
// and optional working confirmation.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - statusCode: HTTP response status code.
//   - encPath: encrypted file path for the working message.
func TestResult(cmd *cobra.Command, statusCode int, encPath string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplTestResult, statusCode, http.StatusText(statusCode))
	if statusCode >= 200 && statusCode < 300 {
		io.sprintf(cmd, config.tplTestWorking, encPath)
	}
}
