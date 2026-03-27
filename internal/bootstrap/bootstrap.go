//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
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
	embedText "github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// registration pairs a command constructor with its group ID.
// An empty groupID marks the command as hidden.
type registration struct {
	cmd     func() *cobra.Command
	groupID string
}

// Initialize registers all ctx subcommands with the root command and
// organizes them into Cobra command groups for structured help output.
//
// Parameters:
//   - cmd: The root cobra command to attach subcommands to
//
// Returns:
//   - *cobra.Command: The same command with all subcommands registered
func Initialize(cmd *cobra.Command) *cobra.Command {
	// Register groups — order determines display order in help.
	cmd.AddGroup(
		&cobra.Group{ID: embedCmd.GroupGettingStarted, Title: desc.Text(embedText.DescKeyGroupGettingStarted)},
		&cobra.Group{ID: embedCmd.GroupContext, Title: desc.Text(embedText.DescKeyGroupContext)},
		&cobra.Group{ID: embedCmd.GroupArtifacts, Title: desc.Text(embedText.DescKeyGroupArtifacts)},
		&cobra.Group{ID: embedCmd.GroupSessions, Title: desc.Text(embedText.DescKeyGroupSessions)},
		&cobra.Group{ID: embedCmd.GroupRuntime, Title: desc.Text(embedText.DescKeyGroupRuntime)},
		&cobra.Group{ID: embedCmd.GroupIntegration, Title: desc.Text(embedText.DescKeyGroupIntegration)},
		&cobra.Group{ID: embedCmd.GroupDiagnostics, Title: desc.Text(embedText.DescKeyGroupDiagnostics)},
		&cobra.Group{ID: embedCmd.GroupUtilities, Title: desc.Text(embedText.DescKeyGroupUtilities)},
	)

	// Assign built-in commands to groups.
	cmd.SetHelpCommandGroupID(embedCmd.GroupUtilities)
	cmd.SetCompletionCommandGroupID(embedCmd.GroupUtilities)

	// Command-to-group mapping. Centralized here so the taxonomy
	// is visible in one place.
	var commands []registration
	commands = append(commands, gettingStarted()...)
	commands = append(commands, contextCmds()...)
	commands = append(commands, artifacts()...)
	commands = append(commands, sessions()...)
	commands = append(commands, runtimeCmds()...)
	commands = append(commands, integrations()...)
	commands = append(commands, diagnostics()...)
	commands = append(commands, utilities()...)
	commands = append(commands, hiddenCmds()...)

	for _, reg := range commands {
		c := reg.cmd()
		if reg.groupID != "" {
			c.GroupID = reg.groupID
		} else {
			c.Hidden = true
		}
		cmd.AddCommand(c)
	}

	return cmd
}

func gettingStarted() []registration {
	return []registration{
		{initialize.Cmd, embedCmd.GroupGettingStarted},
		{status.Cmd, embedCmd.GroupGettingStarted},
		{guide.Cmd, embedCmd.GroupGettingStarted},
	}
}

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

func artifacts() []registration {
	return []registration{
		{decision.Cmd, embedCmd.GroupArtifacts},
		{learning.Cmd, embedCmd.GroupArtifacts},
		{task.Cmd, embedCmd.GroupArtifacts},
	}
}

func sessions() []registration {
	return []registration{
		{journal.Cmd, embedCmd.GroupSessions},
		{memory.Cmd, embedCmd.GroupSessions},
		{remind.Cmd, embedCmd.GroupSessions},
		{pad.Cmd, embedCmd.GroupSessions},
	}
}

func runtimeCmds() []registration {
	return []registration{
		{config.Cmd, embedCmd.GroupRuntime},
		{permission.Cmd, embedCmd.GroupRuntime},
		{pause.Cmd, embedCmd.GroupRuntime},
		{resume.Cmd, embedCmd.GroupRuntime},
	}
}

func integrations() []registration {
	return []registration{
		{hook.Cmd, embedCmd.GroupIntegration},
		{mcp.Cmd, embedCmd.GroupIntegration},
		{watch.Cmd, embedCmd.GroupIntegration},
		{notify.Cmd, embedCmd.GroupIntegration},
		{loop.Cmd, embedCmd.GroupIntegration},
	}
}

func diagnostics() []registration {
	return []registration{
		{doctor.Cmd, embedCmd.GroupDiagnostics},
		{change.Cmd, embedCmd.GroupDiagnostics},
		{dep.Cmd, embedCmd.GroupDiagnostics},
		{why.Cmd, embedCmd.GroupDiagnostics},
	}
}

func utilities() []registration {
	return []registration{
		{reindex.Cmd, embedCmd.GroupUtilities},
	}
}

func hiddenCmds() []registration {
	return []registration{
		{serve.Cmd, ""},
		{site.Cmd, ""},
		{system.Cmd, ""},
	}
}
