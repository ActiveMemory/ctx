//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permission

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/permission/cmd/restore"
	"github.com/ActiveMemory/ctx/internal/cli/permission/cmd/snapshot"
)

// Cmd returns the permission command with subcommands.
//
// The permission command provides utilities for managing Claude Code
// permission snapshots:
//   - snapshot: Save settings.local.json as a golden image
//   - restore: Reset settings.local.json from the golden image
//
// Returns:
//   - *cobra.Command: Configured permission command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyPermission)

	cmd := &cobra.Command{
		Use:   cmd.UsePermission,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(snapshot.Cmd())
	cmd.AddCommand(restore.Cmd())

	return cmd
}
