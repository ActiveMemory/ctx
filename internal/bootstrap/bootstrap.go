//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	embedCmd "github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedText "github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Initialize registers all ctx subcommands with the root command and
// organizes them into Cobra command groups for structured help output.
//
// Parameters:
//   - cmd: The root cobra command to attach subcommands to
//
// Returns:
//   - *cobra.Command: The same command with all subcommands registered
func Initialize(cmd *cobra.Command) *cobra.Command {
	// Register groups: order determines display order in help.
	cmd.AddGroup(
		&cobra.Group{
			ID:    embedCmd.GroupGettingStarted,
			Title: desc.Text(embedText.DescKeyGroupGettingStarted),
		},
		&cobra.Group{
			ID:    embedCmd.GroupContext,
			Title: desc.Text(embedText.DescKeyGroupContext),
		},
		&cobra.Group{
			ID:    embedCmd.GroupArtifacts,
			Title: desc.Text(embedText.DescKeyGroupArtifacts),
		},
		&cobra.Group{
			ID:    embedCmd.GroupSessions,
			Title: desc.Text(embedText.DescKeyGroupSessions),
		},
		&cobra.Group{
			ID:    embedCmd.GroupRuntime,
			Title: desc.Text(embedText.DescKeyGroupRuntime),
		},
		&cobra.Group{
			ID:    embedCmd.GroupIntegration,
			Title: desc.Text(embedText.DescKeyGroupIntegration),
		},
		&cobra.Group{
			ID:    embedCmd.GroupDiagnostics,
			Title: desc.Text(embedText.DescKeyGroupDiagnostics),
		},
		&cobra.Group{
			ID:    embedCmd.GroupShell,
			Title: desc.Text(embedText.DescKeyGroupShell),
		},
	)

	// Route Cobra's built-in help and completion into the
	// narrow Shell group. These are the only commands that
	// live there — everything else has a domain group.
	cmd.SetHelpCommandGroupID(embedCmd.GroupShell)
	cmd.SetCompletionCommandGroupID(embedCmd.GroupShell)

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
