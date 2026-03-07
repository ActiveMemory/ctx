//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"github.com/spf13/cobra"
)

// Cmd returns the pad resolve subcommand.
//
// Returns:
//   - *cobra.Command: Configured resolve subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resolve",
		Short: "Show both sides of a merge conflict",
		Long: `Decrypt and display both sides of a merge conflict for the scratchpad.

Git stores conflict versions as .context/scratchpad.enc.ours and
.context/scratchpad.enc.theirs during a merge conflict. This command
decrypts both and displays them for manual resolution.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResolve(cmd)
		},
	}
}
