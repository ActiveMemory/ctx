//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system" parent command.
//
// Visible subcommands:
//   - resources: Display system resource usage with threshold alerts
//
// Hidden plumbing subcommands (used by skills and automation):
//   - mark-journal: Update journal processing state
//
// Hidden hook subcommands implement Claude Code hook logic as native Go
// binaries and are not intended for direct user invocation.
//
// Returns:
//   - *cobra.Command: Parent command with resource display, plumbing, and hook subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System diagnostics and hook commands",
		Long: `System diagnostics and hook commands.

Subcommands:
  resources            Show system resource usage (memory, swap, disk, load)
  bootstrap            Print context location for AI agents

Plumbing subcommands (used by skills and automation):
  mark-journal         Update journal processing state

Hook subcommands (Claude Code plugin — safe to run manually):
  check-context-size   Context size checkpoint
  check-ceremonies     Session ceremony adoption nudge
  check-persistence    Context persistence nudge
  check-journal        Journal maintenance reminder
  check-resources      Resource pressure warning (DANGER only)
  check-knowledge      Knowledge file growth nudge
  check-version        Version update nudge
  block-non-path-ctx   Block non-PATH ctx invocations
  post-commit          Post-commit context capture nudge
  cleanup-tmp          Remove stale temp files
  qa-reminder          QA reminder before completion`,
	}

	cmd.AddCommand(
		resourcesCmd(),
		bootstrapCmd(),
		markJournalCmd(),
		checkContextSizeCmd(),
		checkPersistenceCmd(),
		checkJournalCmd(),
		checkCeremoniesCmd(),
		checkVersionCmd(),
		blockNonPathCtxCmd(),
		postCommitCmd(),
		cleanupTmpCmd(),
		qaReminderCmd(),
		checkResourcesCmd(),
		checkKnowledgeCmd(),
	)

	return cmd
}
