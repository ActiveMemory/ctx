//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permission

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/permission/cmd/restore"
	"github.com/ActiveMemory/ctx/internal/cli/permission/cmd/snapshot"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
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
	short, long := desc.Command(cmd.DescKeyPermission)

	c := &cobra.Command{
		Use:   cmd.UsePermission,
		Short: short,
		Long:  long,
	}

	c.AddCommand(snapshot.Cmd())
	c.AddCommand(restore.Cmd())

	return c
}
