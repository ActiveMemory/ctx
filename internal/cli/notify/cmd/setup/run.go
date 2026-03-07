//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	notifylib "github.com/ActiveMemory/ctx/internal/notify"
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
	cmd.Print("Enter webhook URL: ")

	scanner := bufio.NewScanner(stdin)
	if !scanner.Scan() {
		return fmt.Errorf("no input received")
	}
	url := strings.TrimSpace(scanner.Text())
	if url == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	if saveErr := notifylib.SaveWebhook(url); saveErr != nil {
		return fmt.Errorf("save webhook: %w", saveErr)
	}

	masked := MaskURL(url)
	cmd.Println("Webhook configured: " + masked)
	cmd.Println("Encrypted at: .context/.notify.enc")

	return nil
}

// MaskURL shows the scheme + host and masks everything after.
//
// Exported for testability.
//
// Parameters:
//   - url: Full webhook URL
//
// Returns:
//   - string: Masked URL safe for display
func MaskURL(url string) string {
	// Find the third slash (end of scheme://host)
	count := 0
	for i, c := range url {
		if c == '/' {
			count++
			if count == 3 {
				return url[:i] + "/***"
			}
		}
	}
	// No path — show as-is but with masked end
	if len(url) > 20 {
		return url[:20] + "***"
	}
	return url
}
