//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system backup" subcommand.
//
// Returns:
//   - *cobra.Command: Configured backup subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemBackup)

	cmd := &cobra.Command{
		Use:   "backup",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().String("scope", archive.BackupScopeAll,
		assets.FlagDesc(assets.FlagDescKeySystemBackupScope),
	)
	cmd.Flags().Bool("json", false,
		assets.FlagDesc(assets.FlagDescKeySystemBackupJson),
	)

	return cmd
}
