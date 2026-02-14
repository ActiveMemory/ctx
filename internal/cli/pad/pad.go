//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd returns the pad command with subcommands.
//
// When invoked without a subcommand, it lists all scratchpad entries.
//
// Returns:
//   - *cobra.Command: Configured pad command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pad",
		Short: "Encrypted scratchpad for sensitive one-liners",
		Long: `Manage an encrypted scratchpad stored in .context/.

Entries are short one-liners encrypted with AES-256-GCM. The key is
stored at .context/.scratchpad.key (gitignored). The encrypted file
(.context/scratchpad.enc) is committed to git.

When invoked without a subcommand, lists all entries.

Subcommands:
  add      Append a new entry
  rm       Remove an entry by number
  edit     Replace an entry by number
  mv       Move an entry to a different position
  resolve  Show both sides of a merge conflict`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd)
		},
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(rmCmd())
	cmd.AddCommand(editCmd())
	cmd.AddCommand(mvCmd())
	cmd.AddCommand(resolveCmd())

	return cmd
}

// runList prints all scratchpad entries numbered 1-based.
func runList(cmd *cobra.Command) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		cmd.Println(msgEmpty)
		return nil
	}

	for i, entry := range entries {
		cmd.Printf("  %d. %s\n", i+1, entry)
	}

	return nil
}

// validateIndex checks that n is a valid 1-based index into entries.
func validateIndex(n int, entries []string) error {
	if n < 1 || n > len(entries) {
		return fmt.Errorf("%s", errEntryRange(n, len(entries)))
	}
	return nil
}
