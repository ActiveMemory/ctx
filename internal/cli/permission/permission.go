//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permission

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/permission/cmd/restore"
	"github.com/ActiveMemory/ctx/internal/cli/permission/cmd/snapshot"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the permission command with subcommands.
//
// The permission command provides utilities for managing Claude
// Code permission snapshots:
//   - snapshot: Save settings.local.json as a golden image
//   - restore: Reset from the golden image
//
// Returns:
//   - *cobra.Command: Configured permission command
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyPermission, cmd.UsePermission,
		snapshot.Cmd(),
		restore.Cmd(),
	)
}
