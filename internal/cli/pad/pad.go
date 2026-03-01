//   /    ctx:                         https://ctx.ist
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
stored at .context/.ctx.key (gitignored). The encrypted file
(.context/scratchpad.enc) is committed to git.

File blobs can be stored as entries using "add --file". Blob entries use
the format "label:::base64data" and are shown as "label [BLOB]" in the
list view. Use "show N" to decode or "show N --out file" to write to disk.

When invoked without a subcommand, lists all entries.

Subcommands:
  show     Output raw text of an entry by number
  add      Append a new entry
  rm       Remove an entry by number
  edit     Replace an entry by number
  mv       Move an entry to a different position
  resolve  Show both sides of a merge conflict
  import   Bulk-import lines from a file
  export   Export blob entries to a directory as files
  merge    Merge entries from scratchpad files`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd)
		},
	}

	cmd.AddCommand(showCmd())
	cmd.AddCommand(addCmd())
	cmd.AddCommand(rmCmd())
	cmd.AddCommand(editCmd())
	cmd.AddCommand(mvCmd())
	cmd.AddCommand(resolveCmd())
	cmd.AddCommand(importCmd())
	cmd.AddCommand(exportCmd())
	cmd.AddCommand(mergeCmd())

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
		cmd.Printf("  %d. %s\n", i+1, displayEntry(entry))
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
