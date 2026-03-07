//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex provides the "ctx learnings reindex" subcommand.
package reindex

import "github.com/spf13/cobra"

// Cmd returns the reindex subcommand for learnings.
//
// Returns:
//   - *cobra.Command: Command for regenerating the LEARNINGS.md index
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reindex",
		Short: "Regenerate the quick-reference index",
		Long: `Regenerate the quick-reference index at the top of LEARNINGS.md.

The index is a compact table showing date and title for each learning,
allowing AI agents to quickly scan entries without reading the full file.

This command is useful after manual edits to LEARNINGS.md or when
migrating existing files to use the index format.

Examples:
  ctx learnings reindex`,
		RunE: run,
	}
}
