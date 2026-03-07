//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"github.com/spf13/cobra"
)

// Cmd returns the pad merge subcommand.
//
// Returns:
//   - *cobra.Command: Configured merge subcommand
func Cmd() *cobra.Command {
	var keyFile string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "merge FILE...",
		Short: "Merge entries from scratchpad files into the current pad",
		Long: `Merge entries from one or more scratchpad files into the current pad.

Each input file is auto-detected as encrypted or plaintext: decryption is
attempted first, and on failure the file is parsed as plain text. Entries
are deduplicated by exact content — position does not matter.

Use --key to provide a key file for encrypted pads from other projects.

Examples:
  ctx pad merge worktree/.context/scratchpad.enc
  ctx pad merge notes.md backup.enc
  ctx pad merge --key /other/.ctx.key foreign.enc
  ctx pad merge --dry-run pad-a.enc pad-b.md`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(cmd, args, keyFile, dryRun)
		},
	}

	cmd.Flags().StringVarP(&keyFile, "key", "k", "",
		"path to key file for decrypting input files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"print what would be merged without writing")

	return cmd
}
