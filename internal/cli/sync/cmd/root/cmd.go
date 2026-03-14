//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx sync" command for reconciling context with codebase.
//
// The command scans the codebase for changes that should be reflected in
// context files, such as new directories, package manager files, and
// configuration files.
//
// Flags:
//   - --dry-run: Show what would change without modifying files
//
// Returns:
//   - *cobra.Command: Configured sync command with flags registered
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := assets.CommandDesc(assets.CmdDescKeySync)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun,
		"dry-run", false, assets.FlagDesc(assets.FlagDescKeySyncDryRun),
	)

	return cmd
}
