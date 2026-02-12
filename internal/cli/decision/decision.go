//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decision provides commands for managing DECISIONS.md.
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
