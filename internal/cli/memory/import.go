//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"github.com/spf13/cobra"
)

// importCmd returns the memory import subcommand.
//
// Returns:
//   - *cobra.Command: command for importing MEMORY.md entries into .context/ files.
func importCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import entries from MEMORY.md into .context/ files",
		Long: `Classify and promote entries from Claude Code's MEMORY.md into
structured .context/ files using heuristic keyword matching.

Each entry is classified as a convention, decision, learning, task,
or skipped (session notes, generic text). Deduplication prevents
re-importing the same entry.

Exit codes:
  0  Imported successfully (or nothing new to import)
  1  MEMORY.md not found`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runImport(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false,
		"Show classification plan without writing",
	)

	return cmd
}
