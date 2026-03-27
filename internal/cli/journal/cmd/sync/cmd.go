//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx journal sync" subcommand.
//
// Scans journal markdowns and syncs their frontmatter lock state into
// .state.json. This is the inverse of "ctx journal lock": the frontmatter
// is treated as the source of truth, and the state is updated to match.
//
// Returns:
//   - *cobra.Command: Command for syncing lock state from frontmatter
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyRecallSync)

	c := &cobra.Command{
		Use:   cmd.UseRecallSync,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	return c
}
