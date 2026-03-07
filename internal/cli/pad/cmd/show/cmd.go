//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// Cmd returns the pad show subcommand.
//
// Outputs the raw text of entry N (1-based) with no numbering prefix.
// Designed for pipe composability:
//
//	ctx pad edit 1 --append "$(ctx pad show 3)"
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func Cmd() *cobra.Command {
	var outPath string

	cmd := &cobra.Command{
		Use:   "show N",
		Short: "Output raw text of an entry by number",
		Long: `Output the raw text of entry N with no numbering prefix.

Designed for unix pipe composability. The output contains just the entry
text followed by a single trailing newline.

For blob entries, the decoded file content is printed (or written to disk
with --out).

Examples:
  ctx pad show 3
  ctx pad show 3 --out ./recovered.md
  ctx pad edit 1 --append "$(ctx pad show 3)"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid index: %s", args[0])
			}
			return runShow(cmd, n, outPath)
		},
	}

	cmd.Flags().StringVar(&outPath, "out", "", "write blob content to a file")

	return cmd
}
