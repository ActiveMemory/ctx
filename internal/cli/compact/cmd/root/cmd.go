//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx compact" command for cleaning up context files.
//
// The command moves completed tasks to a "Completed (Recent)" section,
// optionally archives old content, and removes empty sections from all
// context files.
//
// Flags:
//   - --archive: Create .context/archive/ for old completed tasks
//
// Returns:
//   - *cobra.Command: Configured compact command with flags registered
func Cmd() *cobra.Command {
	var archive bool

	short, long := desc.Command(cmd.DescKeyCompact)

	c := &cobra.Command{
		Use:     cmd.UseCompact,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyCompact),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, archive)
		},
	}

	flagbind.BoolFlag(
		c, &archive,
		cFlag.Archive, flag.DescKeyCompactArchive,
	)

	return c
}
