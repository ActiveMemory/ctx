//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"bufio"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/spf13/cobra"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	notifylib "github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run prompts for a webhook URL and saves it encrypted.
//
// Exported for testability (tests inject a mock stdin).
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: Input source (os.Stdin in production, temp file in tests)
//
// Returns:
//   - error: Non-nil on empty input or save failure
func Run(cmd *cobra.Command, stdin *os.File) error {
	write.SetupPrompt(cmd)

	scanner := bufio.NewScanner(stdin)
	if !scanner.Scan() {
		return ctxerr.NoInput()
	}
	url := strings.TrimSpace(scanner.Text())
	if url == "" {
		return ctxerr.WebhookEmpty()
	}

	if saveErr := notifylib.SaveWebhook(url); saveErr != nil {
		return ctxerr.SaveWebhook(saveErr)
	}

	write.SetupDone(cmd, notifylib.MaskURL(url), crypto.NotifyEnc)

	return nil
}
