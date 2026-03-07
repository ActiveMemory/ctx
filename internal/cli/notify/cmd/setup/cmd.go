//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx notify setup" subcommand.
//
// Returns:
//   - *cobra.Command: Configured setup subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Configure webhook URL",
		Long: `Prompts for a webhook URL and encrypts it using the scratchpad key.

The URL is stored in .context/.notify.enc (encrypted, safe to commit).
The key lives at ~/.ctx/.ctx.key (user-level, never committed).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunSetup(cmd, os.Stdin)
		},
	}
}
