//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decision provides commands for managing DECISIONS.md.
package decision

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/decision/cmd/reindex"
)

// Cmd returns the decisions command with subcommands.
//
// The decisions command provides utilities for managing the DECISIONS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The decisions command with subcommands
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyDecision)
	cmd := &cobra.Command{
		Use:   "decisions",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(reindex.Cmd())

	return cmd
}
