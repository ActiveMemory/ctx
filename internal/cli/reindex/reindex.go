//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"github.com/spf13/cobra"
)

// Cmd returns the reindex convenience command.
//
// The reindex command regenerates the quick-reference index at the top of
// both DECISIONS.md and LEARNINGS.md in a single invocation.
//
// Returns:
//   - *cobra.Command: The reindex command
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reindex",
		Short: "Regenerate indices for DECISIONS.md and LEARNINGS.md",
		Long: `Regenerate the quick-reference index at the top of both DECISIONS.md
and LEARNINGS.md in a single invocation.

This is a convenience wrapper around:
  ctx decisions reindex
  ctx learnings reindex

The index is a compact table showing date and title for each entry,
allowing AI agents to quickly scan entries without reading the full file.

Run this after manual edits to either file or when migrating existing
files to use the index format.

Examples:
  ctx reindex`,
		RunE: runReindex,
	}
}
