//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system backup" subcommand.
//
// Returns:
//   - *cobra.Command: Configured backup subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup context and Claude data",
		Long: `Create timestamped tar.gz archives of project context and/or global
Claude Code data. Optionally copies archives to an SMB share.

Scopes:
  project  .context/, .claude/, ideas/, ~/.bashrc
  global   ~/.claude/ (excludes todos/)
  all      Both project and global (default)

Environment:
  CTX_BACKUP_SMB_URL    - SMB share URL (e.g. smb://host/share)
  CTX_BACKUP_SMB_SUBDIR - Subdirectory on share (default: ctx-sessions)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBackup(cmd)
		},
	}
	cmd.Flags().String("scope", scopeAll, assets.FlagDesc("system.backup.scope"))
	cmd.Flags().Bool("json", false, assets.FlagDesc("system.backup.json"))
	return cmd
}
