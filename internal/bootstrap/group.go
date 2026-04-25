//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/ActiveMemory/ctx/internal/cli/activate"
	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/agent"
	"github.com/ActiveMemory/ctx/internal/cli/change"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/cli/config"
	"github.com/ActiveMemory/ctx/internal/cli/connection"
	"github.com/ActiveMemory/ctx/internal/cli/deactivate"
	"github.com/ActiveMemory/ctx/internal/cli/decision"
	"github.com/ActiveMemory/ctx/internal/cli/doctor"
	"github.com/ActiveMemory/ctx/internal/cli/drift"
	ctxFmt "github.com/ActiveMemory/ctx/internal/cli/fmt"
	"github.com/ActiveMemory/ctx/internal/cli/guide"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	cliHub "github.com/ActiveMemory/ctx/internal/cli/hub"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/cli/journal"
	"github.com/ActiveMemory/ctx/internal/cli/learning"
	"github.com/ActiveMemory/ctx/internal/cli/load"
	"github.com/ActiveMemory/ctx/internal/cli/loop"
	"github.com/ActiveMemory/ctx/internal/cli/mcp"
	"github.com/ActiveMemory/ctx/internal/cli/memory"
	"github.com/ActiveMemory/ctx/internal/cli/pad"
	"github.com/ActiveMemory/ctx/internal/cli/permission"
	"github.com/ActiveMemory/ctx/internal/cli/prune"
	"github.com/ActiveMemory/ctx/internal/cli/reindex"
	"github.com/ActiveMemory/ctx/internal/cli/remind"
	"github.com/ActiveMemory/ctx/internal/cli/serve"
	"github.com/ActiveMemory/ctx/internal/cli/setup"
	"github.com/ActiveMemory/ctx/internal/cli/site"
	"github.com/ActiveMemory/ctx/internal/cli/skill"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/steering"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/sysinfo"
	"github.com/ActiveMemory/ctx/internal/cli/system"
	"github.com/ActiveMemory/ctx/internal/cli/task"
	"github.com/ActiveMemory/ctx/internal/cli/trace"
	"github.com/ActiveMemory/ctx/internal/cli/trigger"
	"github.com/ActiveMemory/ctx/internal/cli/usage"
	"github.com/ActiveMemory/ctx/internal/cli/watch"
	"github.com/ActiveMemory/ctx/internal/cli/why"
	embedCmd "github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// gettingStarted returns command registrations for the getting-started group.
//
// Returns:
//   - []registration: Init, status, and guide commands
func gettingStarted() []registration {
	return []registration{
		{initialize.Cmd, embedCmd.GroupGettingStarted},
		{activate.Cmd, embedCmd.GroupGettingStarted},
		{deactivate.Cmd, embedCmd.GroupGettingStarted},
		{status.Cmd, embedCmd.GroupGettingStarted},
		{guide.Cmd, embedCmd.GroupGettingStarted},
	}
}

// contextCmds returns command registrations for the context
// management group.
//
// These commands operate on the full set of context source-of-truth
// files (TASKS.md, DECISIONS.md, LEARNINGS.md, CONVENTIONS.md):
// adding entries, loading for agents, formatting, reconciling with
// the codebase, detecting drift, and archiving completed work.
//
// Returns:
//   - []registration: Add, load, agent, skill, sync, drift,
//     compact, and fmt commands
func contextCmds() []registration {
	return []registration{
		{add.Cmd, embedCmd.GroupContext},
		{load.Cmd, embedCmd.GroupContext},
		{agent.Cmd, embedCmd.GroupContext},
		{skill.Cmd, embedCmd.GroupContext},
		{sync.Cmd, embedCmd.GroupContext},
		{drift.Cmd, embedCmd.GroupContext},
		{compact.Cmd, embedCmd.GroupContext},
		{ctxFmt.Cmd, embedCmd.GroupContext},
	}
}

// artifacts returns command registrations for the artifacts group.
//
// These commands operate on specific artifact files inside
// .context/: the DECISIONS.md, LEARNINGS.md, and TASKS.md
// stores, plus the `reindex` shortcut that rebuilds the
// decision/learning index tables in a single call.
//
// Returns:
//   - []registration: Decision, learning, task, and reindex
//     commands
func artifacts() []registration {
	return []registration{
		{decision.Cmd, embedCmd.GroupArtifacts},
		{learning.Cmd, embedCmd.GroupArtifacts},
		{task.Cmd, embedCmd.GroupArtifacts},
		{reindex.Cmd, embedCmd.GroupArtifacts},
	}
}

// sessions returns command registrations for the sessions group.
//
// Returns:
//   - []registration: Journal, memory, remind, and pad commands
func sessions() []registration {
	return []registration{
		{journal.Cmd, embedCmd.GroupSessions},
		{memory.Cmd, embedCmd.GroupSessions},
		{remind.Cmd, embedCmd.GroupSessions},
		{pad.Cmd, embedCmd.GroupSessions},
	}
}

// runtimeCmds returns command registrations for the
// runtime configuration group.
//
// Returns:
//   - []registration: Config, permission, hook, and prune commands
func runtimeCmds() []registration {
	return []registration{
		{config.Cmd, embedCmd.GroupRuntime},
		{permission.Cmd, embedCmd.GroupRuntime},
		{hook.Cmd, embedCmd.GroupRuntime},
		{prune.Cmd, embedCmd.GroupRuntime},
	}
}

// integrations returns command registrations for the integrations group.
//
// This group covers commands that connect ctx to external
// systems: AI-tool setup, the ctx Hub server and its clients,
// the MCP server, webhooks, watchers, and loop harnesses.
//
// Returns:
//   - []registration: Setup, steering, trigger, serve, hub,
//     connect, mcp, watch, and loop commands
func integrations() []registration {
	return []registration{
		{setup.Cmd, embedCmd.GroupIntegration},
		{steering.Cmd, embedCmd.GroupIntegration},
		{trigger.Cmd, embedCmd.GroupIntegration},
		{serve.Cmd, embedCmd.GroupIntegration},
		{cliHub.Cmd, embedCmd.GroupIntegration},
		{connection.Cmd, embedCmd.GroupIntegration},
		{mcp.Cmd, embedCmd.GroupIntegration},
		{watch.Cmd, embedCmd.GroupIntegration},
		{loop.Cmd, embedCmd.GroupIntegration},
	}
}

// diagnostics returns command registrations for the diagnostics group.
//
// Returns:
//   - []registration: Doctor, change, why, trace, sysinfo, and usage commands
func diagnostics() []registration {
	return []registration{
		{doctor.Cmd, embedCmd.GroupDiagnostics},
		{change.Cmd, embedCmd.GroupDiagnostics},
		{why.Cmd, embedCmd.GroupDiagnostics},
		{trace.Cmd, embedCmd.GroupDiagnostics},
		{sysinfo.Cmd, embedCmd.GroupDiagnostics},
		{usage.Cmd, embedCmd.GroupDiagnostics},
	}
}

// hiddenCmds returns command registrations that are intentionally
// kept out of `ctx --help` output.
//
// These are genuinely internal commands, not user-facing
// features: `ctx site` generates the journal site consumed by
// `ctx serve`, and `ctx system` hosts the nudge-hook plumbing
// that ctx itself calls via subprocess.
//
// Returns:
//   - []registration: site and system commands with no group
//     assignment (hidden)
func hiddenCmds() []registration {
	return []registration{
		{site.Cmd, ""},
		{system.Cmd, ""},
	}
}
