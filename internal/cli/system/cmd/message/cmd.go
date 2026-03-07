//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system message" subcommand.
//
// Returns:
//   - *cobra.Command: Configured message subcommand with sub-subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "Manage hook message templates",
		Long: `Manage hook message templates.

Hook messages control what text hooks emit. The hook logic (when to
fire, counting, state tracking) is universal. The messages are opinions
that can be customized per-project.

Subcommands:
  list     Show all hook messages with category and override status
  show     Print the effective message template for a hook/variant
  edit     Copy the embedded default to .context/ for editing
  reset    Delete a user override and revert to embedded default`,
	}

	cmd.AddCommand(
		messageListCmd(),
		messageShowCmd(),
		messageEditCmd(),
		messageResetCmd(),
	)

	return cmd
}
