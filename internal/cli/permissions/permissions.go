//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permissions

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/permissions/cmd/restore"
	"github.com/ActiveMemory/ctx/internal/cli/permissions/cmd/snapshot"
)

// Cmd returns the permissions command with subcommands.
//
// The permissions command provides utilities for managing Claude Code
// permission snapshots:
//   - snapshot: Save settings.local.json as a golden image
//   - restore: Reset settings.local.json from the golden image
//
// Returns:
//   - *cobra.Command: Configured permissions command with subcommands
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyPermissions)

	cmd := &cobra.Command{
		Use:   "permissions",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(snapshot.Cmd())
	cmd.AddCommand(restore.Cmd())

	return cmd
}
