//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/agent"
	"github.com/ActiveMemory/ctx/internal/cli/change"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/cli/config"
	"github.com/ActiveMemory/ctx/internal/cli/decision"
	"github.com/ActiveMemory/ctx/internal/cli/dep"
	"github.com/ActiveMemory/ctx/internal/cli/doctor"
	"github.com/ActiveMemory/ctx/internal/cli/drift"
	"github.com/ActiveMemory/ctx/internal/cli/guide"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/cli/journal"
	"github.com/ActiveMemory/ctx/internal/cli/learning"
	"github.com/ActiveMemory/ctx/internal/cli/load"
	"github.com/ActiveMemory/ctx/internal/cli/loop"
	"github.com/ActiveMemory/ctx/internal/cli/mcp"
	"github.com/ActiveMemory/ctx/internal/cli/memory"
	"github.com/ActiveMemory/ctx/internal/cli/notify"
	"github.com/ActiveMemory/ctx/internal/cli/pad"
	"github.com/ActiveMemory/ctx/internal/cli/pause"
	"github.com/ActiveMemory/ctx/internal/cli/permission"
	"github.com/ActiveMemory/ctx/internal/cli/reindex"
	"github.com/ActiveMemory/ctx/internal/cli/remind"
	"github.com/ActiveMemory/ctx/internal/cli/resume"
	"github.com/ActiveMemory/ctx/internal/cli/serve"
	"github.com/ActiveMemory/ctx/internal/cli/site"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/system"
	"github.com/ActiveMemory/ctx/internal/cli/task"
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
		{status.Cmd, embedCmd.GroupGettingStarted},
		{guide.Cmd, embedCmd.GroupGettingStarted},
	}
}

// contextCmds returns command registrations for the context management group.
//
// Returns:
//   - []registration: Add, load, agent, sync, drift, and compact commands
func contextCmds() []registration {
	return []registration{
		{add.Cmd, embedCmd.GroupContext},
		{load.Cmd, embedCmd.GroupContext},
		{agent.Cmd, embedCmd.GroupContext},
		{sync.Cmd, embedCmd.GroupContext},
		{drift.Cmd, embedCmd.GroupContext},
		{compact.Cmd, embedCmd.GroupContext},
	}
}

// artifacts returns command registrations for the artifacts group.
//
// Returns:
//   - []registration: Decision, learning, and task commands
func artifacts() []registration {
	return []registration{
		{decision.Cmd, embedCmd.GroupArtifacts},
		{learning.Cmd, embedCmd.GroupArtifacts},
		{task.Cmd, embedCmd.GroupArtifacts},
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
//   - []registration: Config, permission, pause, and resume commands
func runtimeCmds() []registration {
	return []registration{
		{config.Cmd, embedCmd.GroupRuntime},
		{permission.Cmd, embedCmd.GroupRuntime},
		{pause.Cmd, embedCmd.GroupRuntime},
		{resume.Cmd, embedCmd.GroupRuntime},
	}
}

// integrations returns command registrations for the integrations group.
//
// Returns:
//   - []registration: Hook, mcp, watch, notify, and loop commands
func integrations() []registration {
	return []registration{
		{hook.Cmd, embedCmd.GroupIntegration},
		{mcp.Cmd, embedCmd.GroupIntegration},
		{watch.Cmd, embedCmd.GroupIntegration},
		{notify.Cmd, embedCmd.GroupIntegration},
		{loop.Cmd, embedCmd.GroupIntegration},
	}
}

// diagnostics returns command registrations for the diagnostics group.
//
// Returns:
//   - []registration: Doctor, change, dep, and why commands
func diagnostics() []registration {
	return []registration{
		{doctor.Cmd, embedCmd.GroupDiagnostics},
		{change.Cmd, embedCmd.GroupDiagnostics},
		{dep.Cmd, embedCmd.GroupDiagnostics},
		{why.Cmd, embedCmd.GroupDiagnostics},
	}
}

// utilities returns command registrations for the utilities group.
//
// Returns:
//   - []registration: Reindex command
func utilities() []registration {
	return []registration{
		{reindex.Cmd, embedCmd.GroupUtilities},
	}
}

// hiddenCmds returns command registrations that are not shown in help output.
//
// Returns:
//   - []registration: Serve, site, and system commands with no group assignment
func hiddenCmds() []registration {
	return []registration{
		{serve.Cmd, ""},
		{site.Cmd, ""},
		{system.Cmd, ""},
	}
}
