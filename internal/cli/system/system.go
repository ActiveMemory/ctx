//   /    ctx:                         https://ctx.ist
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
//   - message: Manage hook message templates (list/show/edit/reset)
//
// Hidden plumbing subcommands (used by skills and automation):
//   - mark-journal: Update journal processing state
//   - mark-wrapped-up: Suppress checkpoint nudges after wrap-up
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
  backup               Backup context and Claude data
  resources            Show system resource usage (memory, swap, disk, load)
  bootstrap            Print context location for AI agents
  message              Manage hook message templates (list/show/edit/reset)

Plumbing subcommands (used by skills and automation):
  mark-journal         Update journal processing state
  mark-wrapped-up      Suppress checkpoint nudges after wrap-up
  pause                Pause context hooks for this session
  resume               Resume context hooks for this session
  events               Query the local hook event log

Hook subcommands (Claude Code plugin — safe to run manually):
  context-load-gate           Context file read directive (PreToolUse)
  check-context-size          Context size checkpoint
  check-ceremonies            Session ceremony adoption nudge
  check-persistence           Context persistence nudge
  check-journal               Journal maintenance reminder
  check-resources             Resource pressure warning (DANGER only)
  check-knowledge             Knowledge file growth nudge
  check-reminders             Pending reminders relay
  check-version               Version update nudge
  check-map-staleness         Architecture map staleness nudge
  block-non-path-ctx          Block non-PATH ctx invocations
  block-dangerous-commands    Block dangerous command patterns (project-local)
  check-backup-age            Backup staleness check (project-local)
  post-commit                 Post-commit context capture nudge
  cleanup-tmp                 Remove stale temp files
  qa-reminder                 QA reminder before completion
  specs-nudge                 Plan-to-specs directory nudge (PreToolUse)
  heartbeat                   Session heartbeat webhook (no stdout)`,
	}

	cmd.AddCommand(
		backupCmd(),
		resourcesCmd(),
		bootstrapCmd(),
		messageCmd(),
		markJournalCmd(),
		markWrappedUpCmd(),
		pauseCmd(),
		resumeCmd(),
		eventsCmd(),
		contextLoadGateCmd(),
		checkContextSizeCmd(),
		checkPersistenceCmd(),
		checkJournalCmd(),
		checkCeremoniesCmd(),
		checkRemindersCmd(),
		checkVersionCmd(),
		blockNonPathCtxCmd(),
		postCommitCmd(),
		cleanupTmpCmd(),
		qaReminderCmd(),
		checkResourcesCmd(),
		checkKnowledgeCmd(),
		checkMapStalenessCmd(),
		blockDangerousCommandsCmd(),
		checkBackupAgeCmd(),
		specsNudgeCmd(),
		heartbeatCmd(),
	)

	return cmd
}
