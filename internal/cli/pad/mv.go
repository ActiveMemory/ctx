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

// mvCmd returns the pad mv subcommand.
//
// Returns:
//   - *cobra.Command: Configured mv subcommand
func mvCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mv N M",
		Short: "Move an entry from position N to position M",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			m, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			return runMv(cmd, n, m)
		},
	}
}

// runMv moves entry from 1-based position n to 1-based position m.
func runMv(cmd *cobra.Command, n, m int) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if err := validateIndex(n, entries); err != nil {
		return err
	}
	if err := validateIndex(m, entries); err != nil {
		return err
	}

	// Extract the entry at position n
	entry := entries[n-1]
	// Remove it
	entries = append(entries[:n-1], entries[n:]...)
	// Insert at position m (adjust for 0-based)
	idx := m - 1
	entries = append(entries[:idx], append([]string{entry}, entries[idx:]...)...)

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Moved entry %d to %d.\n", n, m)
	return nil
}
