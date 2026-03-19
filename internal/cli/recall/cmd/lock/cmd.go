//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lock

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx recall lock" subcommand.
//
// Protects journal entries from being overwritten by export --regenerate.
// Locked entries are skipped during export regardless of flags.
//
// Returns:
//   - *cobra.Command: Command for locking journal entries
func Cmd() *cobra.Command {
	var all bool

	short, long := desc.CommandDesc(cmd.DescKeyRecallLock)

	cmd := &cobra.Command{
		Use:   "lock <pattern>",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockUnlock(cmd, args, all, true)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false,
		desc.FlagDesc(flag.DescKeyRecallLockAll),
	)

	return cmd
}
