//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unlock

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx recall unlock" subcommand.
//
// Removes lock protection from journal entries, allowing export
// --regenerate to overwrite them again.
//
// Returns:
//   - *cobra.Command: Command for unlocking journal entries
func Cmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "unlock <pattern>",
		Short: "Remove lock protection from journal entries",
		Long: `Unlock journal entries to allow export --regenerate to overwrite them.

The pattern matches against filenames by slug, date, or short ID (same
matching as export). Unlocking a multi-part entry unlocks all parts.

Examples:
  ctx recall unlock 2026-01-21-session-abc12345.md
  ctx recall unlock abc12345
  ctx recall unlock --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnlock(cmd, args, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Unlock all journal entries")

	return cmd
}
