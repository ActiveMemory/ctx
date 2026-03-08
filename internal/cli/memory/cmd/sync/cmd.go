//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the "ctx memory sync" subcommand.
package sync

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the memory sync subcommand.
//
// Returns:
//   - *cobra.Command: command for syncing MEMORY.md to mirror.
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := assets.CommandDesc(assets.CmdDescKeyMemorySync)
	cmd := &cobra.Command{
		Use:   "sync",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false, assets.FlagDesc(assets.FlagDescKeyMemorySyncDryRun),
	)

	return cmd
}
