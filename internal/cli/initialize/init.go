//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	initroot "github.com/ActiveMemory/ctx/internal/cli/initialize/cmd/root"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core"
	"github.com/ActiveMemory/ctx/internal/config"
)

// PluginInstalled reports whether the ctx plugin is registered in
// ~/.claude/plugins/installed_plugins.json.
// Re-exported from the core subpackage for backward compatibility.
var PluginInstalled = core.PluginInstalled

// PluginEnabledGlobally reports whether the ctx plugin is enabled in
// ~/.claude/settings.json.
// Re-exported from the core subpackage for backward compatibility.
var PluginEnabledGlobally = core.PluginEnabledGlobally

// PluginEnabledLocally reports whether the ctx plugin is enabled in
// .claude/settings.local.json in the current project.
// Re-exported from the core subpackage for backward compatibility.
var PluginEnabledLocally = core.PluginEnabledLocally

// Cmd returns the "ctx init" command for initializing a .context/ directory.
//
// The command creates template files for maintaining persistent context
// for AI coding assistants. Files include constitution rules, tasks,
// decisions, learnings, conventions, and architecture documentation.
//
// Flags:
//   - --force, -f: Overwrite existing context files without prompting
//   - --minimal, -m: Only create essential files
//     (TASKS, DECISIONS, CONSTITUTION)
//   - --merge: Auto-merge ctx content into existing CLAUDE.md and PROMPT.md
//   - --ralph: Use autonomous loop templates (no clarifying questions,
//     one-task-per-iteration, completion signals)
//   - --no-plugin-enable: Skip auto-enabling the ctx plugin in
//     ~/.claude/settings.json
//
// Returns:
//   - *cobra.Command: Configured init command with flags registered
func Cmd() *cobra.Command {
	var (
		force          bool
		minimal        bool
		merge          bool
		ralph          bool
		noPluginEnable bool
		caller         string
	)

	short, long := assets.CommandDesc("initialize")
	cmd := &cobra.Command{
		Use:         "init",
		Short:       short,
		Annotations: map[string]string{config.AnnotationSkipInit: "true"},
		Long:        long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initroot.Run(cmd, force, minimal, merge, ralph, noPluginEnable, caller)
		},
	}

	cmd.Flags().BoolVarP(
		&force,
		"force", "f", false, "Overwrite existing context files",
	)
	cmd.Flags().BoolVarP(
		&minimal,
		"minimal", "m", false,
		"Only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md)",
	)
	cmd.Flags().BoolVar(
		&merge, "merge", false,
		"Auto-merge ctx content into existing CLAUDE.md and PROMPT.md",
	)
	cmd.Flags().BoolVar(
		&ralph, "ralph", false,
		"Agent works autonomously without asking questions",
	)
	cmd.Flags().BoolVar(
		&noPluginEnable, "no-plugin-enable", false,
		"Skip auto-enabling the ctx plugin in global Claude Code settings",
	)
	cmd.Flags().StringVar(
		&caller, "caller", "",
		"Identify the calling tool (e.g. vscode) to tailor output",
	)

	return cmd
}
