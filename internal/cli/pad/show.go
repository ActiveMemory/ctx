//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// showCmd returns the pad show subcommand.
//
// Outputs the raw text of entry N (1-based) with no numbering prefix.
// Designed for pipe composability:
//
//	ctx pad edit 1 --append "$(ctx pad show 3)"
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show N",
		Short: "Output raw text of an entry by number",
		Long: `Output the raw text of entry N with no numbering prefix.

Designed for unix pipe composability. The output contains just the entry
text followed by a single trailing newline.

Examples:
  ctx pad show 3
  ctx pad edit 1 --append "$(ctx pad show 3)"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid index: %s", args[0])
			}
			return runShow(cmd, n)
		},
	}
}

// runShow prints the raw text of entry at 1-based position n.
func runShow(cmd *cobra.Command, n int) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("%s", errEntryRange(n, 0))
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	cmd.Println(entries[n-1])
	return nil
}
