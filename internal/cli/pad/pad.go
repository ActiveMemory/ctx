//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/edit"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/export"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/merge"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/mv"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/normalize"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/resolve"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/rm"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/root"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/tags"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/tag"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Cmd returns the pad command with subcommands.
//
// When invoked without a subcommand, it lists all scratchpad entries.
//
// Returns:
//   - *cobra.Command: Configured pad command with subcommands
func Cmd() *cobra.Command {
	var tagFilters []string

	short, long := desc.Command(cmd.DescKeyPad)
	c := &cobra.Command{
		Use:   cmd.UsePad,
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, readErr := store.ReadEntriesWithIDs()
			if readErr != nil {
				return readErr
			}

			if len(entries) == 0 {
				pad.Empty(cmd)
				return nil
			}

			printed := 0
			for _, entry := range entries {
				if len(tagFilters) > 0 &&
					!tag.MatchAll(entry.Content, tagFilters) {
					continue
				}
				pad.EntryList(
					cmd,
					fmt.Sprintf(
						desc.Text(text.DescKeyWritePadListItem),
						entry.ID,
						blob.DisplayEntry(entry.Content),
					),
				)
				printed++
			}

			if printed == 0 {
				pad.Empty(cmd)
			}

			return nil
		},
	}

	flagbind.StringArrayFlagP(c, &tagFilters,
		cFlag.Tag, cFlag.ShortTag,
		embedFlag.DescKeyPadTag,
	)

	c.AddCommand(show.Cmd())
	c.AddCommand(add.Cmd())
	c.AddCommand(rm.Cmd())
	c.AddCommand(edit.Cmd())
	c.AddCommand(mv.Cmd())
	c.AddCommand(resolve.Cmd())
	c.AddCommand(root.Cmd())
	c.AddCommand(export.Cmd())
	c.AddCommand(merge.Cmd())
	c.AddCommand(normalize.Cmd())
	c.AddCommand(tags.Cmd())

	return c
}
