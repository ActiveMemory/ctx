//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx permissions restore" subcommand.
//
// Returns:
//   - *cobra.Command: Configured restore subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Reset settings.local.json from golden image",
		Long: `Replace .claude/settings.local.json with the golden image.

Prints a diff of dropped (session-accumulated) and restored permissions.
No-op if the files already match.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunRestore(cmd)
		},
	}
}
