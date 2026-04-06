//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lock

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx journal lock" subcommand.
//
// Protects journal entries from being overwritten by export --regenerate.
// Locked entries are skipped during export regardless of flags.
//
// Returns:
//   - *cobra.Command: Command for locking journal entries
func Cmd() *cobra.Command {
	var all bool

	short, long := desc.Command(cmd.DescKeyJournalLock)

	c := &cobra.Command{
		Use:     cmd.UseJournalLock,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalLock),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, all)
		},
	}

	flagbind.BoolFlag(
		c, &all,
		cFlag.All, flag.DescKeyJournalLockAll,
	)

	return c
}
