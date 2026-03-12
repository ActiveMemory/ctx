//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"net/http"

	"github.com/ActiveMemory/ctx/internal/write/config"
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
	cmd.Print(config.TplSetupPrompt)
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
	cmd.Println(fmt.Sprintf(config.TplSetupDone, maskedURL, encPath))
}

// TestNoWebhook prints the message when no webhook is configured.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func TestNoWebhook(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplTestNoWebhook)
}

// TestFiltered prints the notice when the test event is filtered.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func TestFiltered(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplTestFiltered)
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
	cmd.Println(fmt.Sprintf(config.TplTestResult, statusCode, http.StatusText(statusCode)))
	if statusCode >= 200 && statusCode < 300 {
		cmd.Println(fmt.Sprintf(config.TplTestWorking, encPath))
	}
}
