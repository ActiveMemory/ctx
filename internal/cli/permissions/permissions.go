//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permissions

import (
	"github.com/spf13/cobra"

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
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage permission snapshots",
		Long: `Manage Claude Code permission snapshots.

Save a curated settings.local.json as a golden image, then restore
at session start to automatically drop session-accumulated permissions.

Subcommands:
  snapshot  Save settings.local.json as golden image
  restore   Reset settings.local.json from golden image`,
	}

	cmd.AddCommand(snapshot.Cmd())
	cmd.AddCommand(restore.Cmd())

	return cmd
}
