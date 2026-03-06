//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory implements the "ctx memory" CLI command group for
// bridging Claude Code's auto memory into .context/.
package memory

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx memory" parent command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Bridge Claude Code auto memory into .context/",
		Long: `Bridge Claude Code's auto memory (MEMORY.md) into .context/.

Discovers MEMORY.md from ~/.claude/projects/, mirrors it into
.context/memory/mirror.md (git-tracked), and detects drift.

Subcommands:
  sync       Copy MEMORY.md to mirror, archive previous version
  status     Show drift, timestamps, and entry counts
  diff       Show what changed since last sync
  import     Classify and promote entries to .context/ files
  publish    Push curated .context/ content to MEMORY.md
  unpublish  Remove published block from MEMORY.md`,
	}

	cmd.AddCommand(
		syncCmd(),
		statusCmd(),
		diffCmd(),
		importCmd(),
		publishCmd(),
		unpublishCmd(),
	)

	return cmd
}
