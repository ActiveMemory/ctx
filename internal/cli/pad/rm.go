//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"strconv"

	"github.com/spf13/cobra"
)

// rmCmd returns the pad rm subcommand.
//
// Returns:
//   - *cobra.Command: Configured rm subcommand
func rmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm N",
		Short: "Remove an entry by number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return runRm(cmd, n)
		},
	}
}

// runRm removes entry at 1-based position n.
func runRm(cmd *cobra.Command, n int) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}

	entries = append(entries[:n-1], entries[n:]...)

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Removed entry %d.\n", n)
	return nil
}
