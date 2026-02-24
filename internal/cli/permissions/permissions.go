//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permissions

import (
	"github.com/spf13/cobra"
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

	cmd.AddCommand(snapshotCmd())
	cmd.AddCommand(restoreCmd())

	return cmd
}

// snapshotCmd returns the permissions snapshot subcommand.
func snapshotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "snapshot",
		Short: "Save settings.local.json as golden image",
		Long: `Save .claude/settings.local.json as the golden image.

The golden file (.claude/settings.golden.json) is a byte-for-byte copy
of the current settings. It is meant to be committed to version control
and shared with the team.

Overwrites any existing golden file.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSnapshot(cmd)
		},
	}
}

// restoreCmd returns the permissions restore subcommand.
func restoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Reset settings.local.json from golden image",
		Long: `Replace .claude/settings.local.json with the golden image.

Prints a diff of dropped (session-accumulated) and restored permissions.
No-op if the files already match.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runRestore(cmd)
		},
	}
}
