//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx compact" command for cleaning up context files.
//
// The command moves completed tasks to a "Completed (Recent)" section,
// optionally archives old content, and removes empty sections from all
// context files.
//
// Flags:
//   - --archive: Create .context/archive/ for old completed tasks
//
// Returns:
//   - *cobra.Command: Configured compact command with flags registered
func Cmd() *cobra.Command {
	var archive bool

	cmd := &cobra.Command{
		Use:   "compact",
		Short: "Archive completed tasks and clean up context",
		Long: `Consolidate and clean up context files.

Actions performed:
  - Move completed tasks to "Completed (Recent)" section
  - Archive old completed tasks (with --archive)
  - Remove empty sections from context files
  - Report on potential duplicates

Use --archive to create .context/archive/ for old content.

Examples:
  ctx compact                  # Clean up context, move completed tasks
  ctx compact --archive        # Also archive old tasks to .context/archive/
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompact(cmd, archive)
		},
	}

	cmd.Flags().BoolVar(
		&archive,
		"archive",
		false,
		"Create .context/archive/ for old content",
	)

	return cmd
}
