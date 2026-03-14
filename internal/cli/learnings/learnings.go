//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learnings

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/learnings/cmd/reindex"
)

// Cmd returns the learnings command with subcommands.
//
// The learnings command provides utilities for managing the LEARNINGS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The learnings command with subcommands
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyLearnings)
	cmd := &cobra.Command{
		Use:   "learnings",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(reindex.Cmd())

	return cmd
}
