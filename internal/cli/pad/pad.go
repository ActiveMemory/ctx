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
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/imp"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/merge"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/mv"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/resolve"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/rm"
	"github.com/ActiveMemory/ctx/internal/cli/pad/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Cmd returns the pad command with subcommands.
//
// When invoked without a subcommand, it lists all scratchpad entries.
//
// Returns:
//   - *cobra.Command: Configured pad command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyPad)
	c := &cobra.Command{
		Use:   cmd.UsePad,
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := core.ReadEntries()
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				pad.Empty(cmd)
				return nil
			}

			for i, entry := range entries {
				pad.EntryList(cmd, fmt.Sprintf(desc.Text(text.DescKeyWritePadListItem), i+1, core.DisplayEntry(entry)))
			}

			return nil
		},
	}

	c.AddCommand(show.Cmd())
	c.AddCommand(add.Cmd())
	c.AddCommand(rm.Cmd())
	c.AddCommand(edit.Cmd())
	c.AddCommand(mv.Cmd())
	c.AddCommand(resolve.Cmd())
	c.AddCommand(imp.Cmd())
	c.AddCommand(export.Cmd())
	c.AddCommand(merge.Cmd())

	return c
}
