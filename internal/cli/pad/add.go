//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"github.com/spf13/cobra"
)

// addCmd returns the pad add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add TEXT",
		Short: "Append a new entry to the scratchpad",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0])
		},
	}
}

// runAdd appends a new entry and prints confirmation.
func runAdd(cmd *cobra.Command, text string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	entries = append(entries, text)

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Added entry %d.\n", len(entries))
	return nil
}
