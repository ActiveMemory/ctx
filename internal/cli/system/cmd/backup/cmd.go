//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system backup" subcommand.
//
// Returns:
//   - *cobra.Command: Configured backup subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemBackup)

	cmd := &cobra.Command{
		Use:   "backup",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().String("scope", archive.BackupScopeAll,
		desc.Flag(flag.DescKeySystemBackupScope),
	)
	cmd.Flags().Bool("json", false,
		desc.Flag(flag.DescKeySystemBackupJson),
	)

	return cmd
}
