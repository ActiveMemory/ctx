//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unlock

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx journal unlock" subcommand.
//
// Removes lock protection from journal entries, allowing export
// --regenerate to overwrite them again.
//
// Returns:
//   - *cobra.Command: Command for unlocking journal entries
func Cmd() *cobra.Command {
	var all bool

	short, long := desc.Command(cmd.DescKeyJournalUnlock)

	c := &cobra.Command{
		Use:     cmd.UseJournalUnlock,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalUnlock),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, all)
		},
	}

	flagbind.BoolFlag(
		c, &all,
		cFlag.All, flag.DescKeyJournalUnlockAll,
	)

	return c
}
