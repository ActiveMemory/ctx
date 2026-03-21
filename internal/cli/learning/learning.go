//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learning

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/learning/cmd/reindex"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the learning command with subcommands.
//
// The learning command provides utilities for managing the LEARNINGS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The learning command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyLearning)
	c := &cobra.Command{
		Use:   cmd.UseLearning,
		Short: short,
		Long:  long,
	}

	c.AddCommand(reindex.Cmd())

	return c
}
