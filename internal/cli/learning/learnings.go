//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learning

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/learning/cmd/reindex"
)

// Cmd returns the learning command with subcommands.
//
// The learning command provides utilities for managing the LEARNINGS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The learning command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyLearning)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyLearning,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(reindex.Cmd())

	return cmd
}
