//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/block_dangerous_command"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/block_non_path_ctx"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_backup_age"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_ceremony"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_context_size"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_freshness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_hub_sync"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_journal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_knowledge"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_map_staleness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_memory_drift"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_persistence"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_reminder"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_resource"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_skill_discovery"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_task_completion"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/check_version"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/context_load_gate"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/heartbeat"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/mark_journal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/mark_wrapped_up"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/pause"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/post_commit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/qa_reminder"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/resume"
	sessEvent "github.com/ActiveMemory/ctx/internal/cli/system/cmd/session_event"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/specs_nudge"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system" parent command.
//
// Hosts hidden Claude Code hook plumbing only. User-facing
// maintenance commands (backup, bootstrap, event, message, prune,
// resource, stats) have been promoted to top-level commands and
// are registered separately in internal/bootstrap/group.go.
//
// Hook subcommands implement Claude Code hook logic as native Go
// binaries and are not intended for direct user invocation.
//
// Returns:
//   - *cobra.Command: Parent command with hook plumbing subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeySystem, cmd.UseSystem,
		block_dangerous_command.Cmd(),
		block_non_path_ctx.Cmd(),
		check_backup_age.Cmd(),
		check_ceremony.Cmd(),
		check_context_size.Cmd(),
		check_freshness.Cmd(),
		check_hub_sync.Cmd(),
		check_journal.Cmd(),
		check_knowledge.Cmd(),
		check_map_staleness.Cmd(),
		check_memory_drift.Cmd(),
		check_persistence.Cmd(),
		check_skill_discovery.Cmd(),
		check_reminder.Cmd(),
		check_resource.Cmd(),
		check_task_completion.Cmd(),
		check_version.Cmd(),
		context_load_gate.Cmd(),
		heartbeat.Cmd(),
		mark_journal.Cmd(),
		mark_wrapped_up.Cmd(),
		pause.Cmd(),
		post_commit.Cmd(),
		qa_reminder.Cmd(),
		resume.Cmd(),
		sessEvent.Cmd(),
		specs_nudge.Cmd(),
	)
}
