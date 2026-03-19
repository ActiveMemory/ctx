//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decision provides commands for managing DECISIONS.md.
package decision

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/decision/cmd/reindex"
)

// Cmd returns the decision command with subcommands.
//
// The decision command provides utilities for managing the DECISIONS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The decision command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyDecision)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyDecision,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(reindex.Cmd())

	return cmd
}
