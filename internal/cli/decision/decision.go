//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decisions provides commands for managing DECISIONS.md.
package decision

import (
	"github.com/spf13/cobra"
)

// Cmd returns the decisions command with subcommands.
//
// The decisions command provides utilities for managing the DECISIONS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The decisions command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decisions",
		Short: "Manage DECISIONS.md file",
		Long: `Manage the DECISIONS.md file and its quick-reference index.

The decisions file maintains an auto-generated index at the top for quick
scanning. Use the subcommands to manage this index.

Subcommands:
  reindex    Regenerate the quick-reference index

Examples:
  ctx decisions reindex`,
	}

	cmd.AddCommand(reindexCmd())

	return cmd
}

// reindexCmd returns the reindex subcommand.
//
// Returns:
//   - *cobra.Command: Command for regenerating the DECISIONS.md index
func reindexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reindex",
		Short: "Regenerate the quick-reference index",
		Long: `Regenerate the quick-reference index at the top of DECISIONS.md.

The index is a compact table showing date and title for each decision,
allowing AI agents to quickly scan entries without reading the full file.

This command is useful after manual edits to DECISIONS.md or when
migrating existing files to use the index format.

Examples:
  ctx decisions reindex`,
		RunE: runReindex,
	}
}
