//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
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
		caller         string
	)

	short, long := desc.Command(cmd.DescKeyInitialize)
	c := &cobra.Command{
		Use:         cmd.UseInit,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, force, minimal, merge, ralph, noPluginEnable, caller)
		},
	}

	c.Flags().BoolVarP(
		&force,
		cFlag.Force, cFlag.ShortForce, false,
		desc.Flag(flag.DescKeyInitializeForce),
	)
	c.Flags().BoolVarP(
		&minimal,
		cFlag.Minimal, cFlag.ShortMinimal, false,
		desc.Flag(flag.DescKeyInitializeMinimal),
	)
	c.Flags().BoolVar(
		&merge, cFlag.Merge, false,
		desc.Flag(flag.DescKeyInitializeMerge),
	)
	c.Flags().BoolVar(
		&ralph, cFlag.Ralph, false,
		desc.Flag(flag.DescKeyInitializeRalph),
	)
	c.Flags().BoolVar(
		&noPluginEnable, cFlag.NoPluginEnable, false,
		desc.Flag(flag.DescKeyInitializeNoPluginEnable),
	)
	c.Flags().StringVar(
		&caller, "caller", "",
		"Identify the calling tool (e.g. vscode) to tailor output",
	)

	return c
}
