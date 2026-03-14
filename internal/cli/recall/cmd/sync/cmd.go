//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx recall sync" subcommand.
//
// Scans journal markdowns and syncs their frontmatter lock state into
// .state.json. This is the inverse of "ctx recall lock": the frontmatter
// is treated as the source of truth, and state is updated to match.
//
// Returns:
//   - *cobra.Command: Command for syncing lock state from frontmatter
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyRecallSync)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	return cmd
}
