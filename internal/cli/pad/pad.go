//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/edit"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/export"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/imp"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/merge"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/mv"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/resolve"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/rm"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
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
stored at ~/.ctx/.ctx.key (global, user-level). The encrypted file
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
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd)
		},
	}

	cmd.AddCommand(show.Cmd())
	cmd.AddCommand(add.Cmd())
	cmd.AddCommand(rm.Cmd())
	cmd.AddCommand(edit.Cmd())
	cmd.AddCommand(mv.Cmd())
	cmd.AddCommand(resolve.Cmd())
	cmd.AddCommand(imp.Cmd())
	cmd.AddCommand(export.Cmd())
	cmd.AddCommand(merge.Cmd())

	return cmd
}

// runList prints all scratchpad entries numbered 1-based.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read failure
func runList(cmd *cobra.Command) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		cmd.Println(core.MsgEmpty)
		return nil
	}

	for i, entry := range entries {
		cmd.Println(fmt.Sprintf("  %d. %s", i+1, core.DisplayEntry(entry)))
	}

	return nil
}
