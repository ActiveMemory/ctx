//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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
	"github.com/ActiveMemory/ctx/internal/write"
)

// Cmd returns the pad command with subcommands.
//
// When invoked without a subcommand, it lists all scratchpad entries.
//
// Returns:
//   - *cobra.Command: Configured pad command with subcommands
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyPad)
	cmd := &cobra.Command{
		Use:   "pad",
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
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
		write.PadEmpty(cmd)
		return nil
	}

	for i, entry := range entries {
		cmd.Println(fmt.Sprintf("  %d. %s", i+1, core.DisplayEntry(entry)))
	}

	return nil
}
