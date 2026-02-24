//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"os"
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

// runShow prints the raw text of entry at 1-based position n.
func runShow(cmd *cobra.Command, n int, outPath string) error {
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

	entry := entries[n-1]

	if label, data, ok := splitBlob(entry); ok {
		_ = label
		if outPath != "" {
			if err := os.WriteFile(outPath, data, 0600); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
			cmd.Printf("Wrote %d bytes to %s\n", len(data), outPath)
			return nil
		}
		cmd.Print(string(data))
		return nil
	}

	// Non-blob entry.
	if outPath != "" {
		return fmt.Errorf("--out can only be used with blob entries")
	}

	cmd.Println(entry)
	return nil
}
