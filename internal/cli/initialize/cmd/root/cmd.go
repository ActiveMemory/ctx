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
	"github.com/ActiveMemory/ctx/internal/flagbind"
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
//   - --merge: Auto-merge ctx content into existing CLAUDE.md
//   - --no-plugin-enable: Skip auto-enabling the ctx plugin in
//     ~/.claude/settings.json
//   - --no-steering-init: Skip scaffolding foundation steering
//     files in .context/steering/
//
// Returns:
//   - *cobra.Command: Configured init command with flags registered
func Cmd() *cobra.Command {
	var (
		force          bool
		minimal        bool
		merge          bool
		noPluginEnable bool
		noSteeringInit bool
		caller         string
	)

	short, long := desc.Command(cmd.DescKeyInitialize)
	c := &cobra.Command{
		Use:         cmd.UseInit,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		Example:     desc.Example(cmd.DescKeyInitialize),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(
				cmd, force, minimal, merge,
				noPluginEnable, noSteeringInit, caller,
			)
		},
	}

	flagbind.BindBoolFlagsP(c,
		[]*bool{&force, &minimal},
		[]string{cFlag.Force, cFlag.Minimal},
		[]string{cFlag.ShortForce, cFlag.ShortMinimal},
		[]string{
			flag.DescKeyInitializeForce,
			flag.DescKeyInitializeMinimal,
		},
	)
	flagbind.BindBoolFlags(c,
		[]*bool{&merge, &noPluginEnable, &noSteeringInit},
		[]string{
			cFlag.Merge, cFlag.NoPluginEnable,
			cFlag.NoSteeringInit,
		},
		[]string{
			flag.DescKeyInitializeMerge,
			flag.DescKeyInitializeNoPluginEnable,
			flag.DescKeyInitializeNoSteeringInit,
		},
	)
	flagbind.StringFlag(c, &caller, cFlag.Caller, flag.DescKeyInitializeCaller)

	return c
}
