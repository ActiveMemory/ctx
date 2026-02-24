//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	notifylib "github.com/ActiveMemory/ctx/internal/notify"
)

// setupCmd returns the "ctx notify setup" subcommand.
func setupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Configure webhook URL",
		Long: `Prompts for a webhook URL and encrypts it using the scratchpad key.

The URL is stored in .context/.notify.enc (encrypted, safe to commit).
The key lives at .context/.context.key (gitignored, never committed).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetup(cmd, os.Stdin)
		},
	}
}

func runSetup(cmd *cobra.Command, stdin *os.File) error {
	cmd.Print("Enter webhook URL: ")

	scanner := bufio.NewScanner(stdin)
	if !scanner.Scan() {
		return fmt.Errorf("no input received")
	}
	url := strings.TrimSpace(scanner.Text())
	if url == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	if err := notifylib.SaveWebhook(url); err != nil {
		return fmt.Errorf("save webhook: %w", err)
	}

	masked := maskURL(url)
	cmd.Println("Webhook configured: " + masked)
	cmd.Println("Encrypted at: .context/.notify.enc")

	return nil
}

// maskURL shows the scheme + host and masks everything after.
func maskURL(url string) string {
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
