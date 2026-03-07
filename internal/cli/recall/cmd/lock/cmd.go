//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lock

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx recall lock" subcommand.
//
// Protects journal entries from being overwritten by export --regenerate.
// Locked entries are skipped during export regardless of flags.
//
// Returns:
//   - *cobra.Command: Command for locking journal entries
func Cmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "lock <pattern>",
		Short: "Protect journal entries from export regeneration",
		Long: `Lock journal entries to prevent export --regenerate from overwriting them.

Locked entries are skipped during export regardless of --regenerate or --force.
Use "ctx recall unlock" to remove the protection.

The pattern matches against filenames by slug, date, or short ID (same
matching as export). Locking a multi-part entry locks all parts.

The lock is recorded in .context/journal/.state.json (source of truth) and
a "locked: true" line is added to the file's YAML frontmatter for visibility.

Examples:
  ctx recall lock 2026-01-21-session-abc12345.md
  ctx recall lock abc12345
  ctx recall lock --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockUnlock(cmd, args, all, true)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Lock all journal entries")

	return cmd
}
