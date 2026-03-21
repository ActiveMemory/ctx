//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/backup"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/block_dangerous_commands"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/block_non_path_ctx"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/bootstrap"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_backup_age"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_ceremonies"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_context_size"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_freshness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_journal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_knowledge"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_map_staleness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_memory_drift"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_persistence"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_reminders"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_resources"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_task_completion"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_version"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/context_load_gate"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/events"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/heartbeat"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/mark_journal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/mark_wrapped_up"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/pause"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/post_commit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/prune"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/qa_reminder"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/resources"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/resume"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/specs_nudge"
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
	short, long := desc.CommandDesc(cmd.DescKeySystem)

	cmd := &cobra.Command{
		Use:   cmd.UseSystem,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(
		backup.Cmd(),
		block_dangerous_commands.Cmd(),
		block_non_path_ctx.Cmd(),
		bootstrap.Cmd(),
		check_backup_age.Cmd(),
		check_ceremonies.Cmd(),
		check_context_size.Cmd(),
		check_freshness.Cmd(),
		check_journal.Cmd(),
		check_knowledge.Cmd(),
		check_map_staleness.Cmd(),
		check_memory_drift.Cmd(),
		check_persistence.Cmd(),
		check_reminders.Cmd(),
		check_resources.Cmd(),
		check_task_completion.Cmd(),
		check_version.Cmd(),
		context_load_gate.Cmd(),
		events.Cmd(),
		heartbeat.Cmd(),
		mark_journal.Cmd(),
		mark_wrapped_up.Cmd(),
		message.Cmd(),
		pause.Cmd(),
		post_commit.Cmd(),
		prune.Cmd(),
		qa_reminder.Cmd(),
		resources.Cmd(),
		resume.Cmd(),
		specs_nudge.Cmd(),
		stats.Cmd(),
	)

	return cmd
}
