//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"github.com/spf13/cobra"
)

// exportCmd returns the pad export subcommand.
//
// Returns:
//   - *cobra.Command: command for exporting blob entries to files.
func exportCmd() *cobra.Command {
	var force, dryRun bool

	cmd := &cobra.Command{
		Use:   "export [DIR]",
		Short: "Export blob entries to a directory as files",
		Long: `Export all blob entries from the scratchpad to a directory as files.
Each blob's label becomes the filename. Non-blob entries are skipped.

When a file already exists, a unix timestamp is prepended to avoid
collisions. Use --force to overwrite instead.

Examples:
  ctx pad export
  ctx pad export ./ideas
  ctx pad export --dry-run
  ctx pad export --force ./backup`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return runExport(cmd, dir, force, dryRun)
		},
	}

	cmd.Flags().BoolVarP(
		&force, "force", "f", false,
		"overwrite existing files instead of timestamping",
	)
	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false,
		"print what would be exported without writing",
	)

	return cmd
}
