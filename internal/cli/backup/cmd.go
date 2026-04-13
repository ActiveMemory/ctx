//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx backup" top-level command.
//
// Returns:
//   - *cobra.Command: Configured backup command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyBackup)

	c := &cobra.Command{
		Use:     cmd.UseBackup,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyBackup),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	c.Flags().String(cFlag.Scope, archive.BackupScopeAll,
		desc.Flag(flag.DescKeyBackupScope),
	)
	c.Flags().Bool(cFlag.JSON, false,
		desc.Flag(flag.DescKeyBackupJson),
	)

	return c
}
