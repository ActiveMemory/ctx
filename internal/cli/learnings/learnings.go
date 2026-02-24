//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learnings

import (
	"github.com/spf13/cobra"
)

// Cmd returns the learnings command with subcommands.
//
// The learnings command provides utilities for managing the LEARNINGS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The learnings command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "learnings",
		Short: "Manage LEARNINGS.md file",
		Long: `Manage the LEARNINGS.md file and its quick-reference index.

The learnings file maintains an auto-generated index at the top for quick
scanning. Use the subcommands to manage this index.

Subcommands:
  reindex    Regenerate the quick-reference index

Examples:
  ctx learnings reindex`,
	}

	cmd.AddCommand(reindexCmd())

	return cmd
}
