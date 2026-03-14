//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

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
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyInitialize)
	cmd := &cobra.Command{
		Use:         "init",
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, force, minimal, merge, ralph, noPluginEnable)
		},
	}

	cmd.Flags().BoolVarP(
		&force,
		"force", "f", false, assets.FlagDesc(assets.FlagDescKeyInitializeForce),
	)
	cmd.Flags().BoolVarP(
		&minimal,
		"minimal", "m", false,
		assets.FlagDesc(assets.FlagDescKeyInitializeMinimal),
	)
	cmd.Flags().BoolVar(
		&merge, "merge", false,
		assets.FlagDesc(assets.FlagDescKeyInitializeMerge),
	)
	cmd.Flags().BoolVar(
		&ralph, "ralph", false,
		assets.FlagDesc(assets.FlagDescKeyInitializeRalph),
	)
	cmd.Flags().BoolVar(
		&noPluginEnable, "no-plugin-enable", false,
		assets.FlagDesc(assets.FlagDescKeyInitializeNoPluginEnable),
	)

	return cmd
}
