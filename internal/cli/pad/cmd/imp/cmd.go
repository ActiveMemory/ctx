//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

import (
	"github.com/spf13/cobra"
)

// Cmd returns the pad import subcommand.
//
// Returns:
//   - *cobra.Command: Configured import subcommand
func Cmd() *cobra.Command {
	var blobs bool

	cmd := &cobra.Command{
		Use:   "import FILE",
		Short: "Bulk-import lines from a file into the scratchpad",
		Long: `Import lines from a file into the scratchpad. Each non-empty line
becomes a separate entry. Use "-" to read from stdin.

With --blobs, import all first-level files from a directory as blob entries.
Each file becomes a blob with the filename as its label. Subdirectories and
non-regular files are skipped.

Examples:
  ctx pad import notes.txt
  grep pattern file | ctx pad import -
  ctx pad import --blobs ./ideas/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if blobs {
				return runImportBlobs(cmd, args[0])
			}
			return runImport(cmd, args[0])
		},
	}

	cmd.Flags().BoolVar(&blobs, "blobs", false,
		"import first-level files from a directory as blob entries")

	return cmd
}
