//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Copy MEMORY.md to mirror, archive previous version",
		Long: `Copy Claude Code's MEMORY.md to .context/memory/mirror.md.

Archives the previous mirror before overwriting. Reports line counts
and drift since last sync.

Exit codes:
  0  Synced successfully
  1  MEMORY.md not found (auto memory not active)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSync(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false, "Show what would happen without writing",
	)

	return cmd
}
