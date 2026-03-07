//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/backup"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/blockdangerouscommands"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/blocknonpathctx"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/bootstrap"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkbackupage"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkceremonies"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkcontextsize"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkjournal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkknowledge"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkmapstaleness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkmemorydrift"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkpersistence"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkreminders"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkresources"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checktaskcompletion"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkversion"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/contextloadgate"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/events"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/heartbeat"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/markjournal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/markwrappedup"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/pause"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/postcommit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/prune"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/qareminder"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/resources"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/resume"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/specsnudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/stats"
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

  stats                Show session token usage stats

Plumbing subcommands (used by skills and automation):
  mark-journal         Update journal processing state
  mark-wrapped-up      Suppress checkpoint nudges after wrap-up
  pause                Pause context hooks for this session
  resume               Resume context hooks for this session
  prune                Clean stale per-session state files
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
  check-task-completion       Task completion nudge after edits
  post-commit                 Post-commit context capture nudge
  qa-reminder                 QA reminder before completion
  specs-nudge                 Plan-to-specs directory nudge (PreToolUse)
  check-memory-drift          Memory drift nudge (MEMORY.md changed)
  heartbeat                   Session heartbeat webhook (no stdout)`,
	}

	cmd.AddCommand(
		backup.Cmd(),
		resources.Cmd(),
		stats.Cmd(),
		bootstrap.Cmd(),
		message.Cmd(),
		markjournal.Cmd(),
		markwrappedup.Cmd(),
		pause.Cmd(),
		resume.Cmd(),
		prune.Cmd(),
		events.Cmd(),
		contextloadgate.Cmd(),
		checkcontextsize.Cmd(),
		checkpersistence.Cmd(),
		checkjournal.Cmd(),
		checkceremonies.Cmd(),
		checkreminders.Cmd(),
		checkversion.Cmd(),
		blocknonpathctx.Cmd(),
		checktaskcompletion.Cmd(),
		postcommit.Cmd(),
		qareminder.Cmd(),
		checkresources.Cmd(),
		checkknowledge.Cmd(),
		checkmapstaleness.Cmd(),
		blockdangerouscommands.Cmd(),
		checkbackupage.Cmd(),
		specsnudge.Cmd(),
		checkmemorydrift.Cmd(),
		heartbeat.Cmd(),
	)

	return cmd
}
