//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx permissions snapshot" subcommand.
//
// Returns:
//   - *cobra.Command: Configured snapshot subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "snapshot",
		Short: "Save settings.local.json as golden image",
		Long: `Save .claude/settings.local.json as the golden image.

The golden file (.claude/settings.golden.json) is a byte-for-byte copy
of the current settings. It is meant to be committed to version control
and shared with the team.

Overwrites any existing golden file.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunSnapshot(cmd)
		},
	}
}
