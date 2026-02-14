//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"strconv"

	"github.com/spf13/cobra"
)

// editCmd returns the pad edit subcommand.
//
// Returns:
//   - *cobra.Command: Configured edit subcommand
func editCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit N TEXT",
		Short: "Replace an entry by number",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return runEdit(cmd, n, args[1])
		},
	}
}

// runEdit replaces entry at 1-based position n with new text.
func runEdit(cmd *cobra.Command, n int, text string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	entries[n-1] = text

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Updated entry %d.\n", n)
	return nil
}
