//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/add/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx add" command for appending entries to context files.
//
// Supported types are defined in [config.FileType] (both singular and plural
// forms accepted, e.g., "decision" or "decisions"). Content can be provided
// via command argument, --file flag, or stdin pipe.
//
// Flags:
//   - --priority, -p: Priority level for tasks (high, medium, low)
//   - --section, -s: Target section within the file
//   - --file, -f: Read content from a file instead of argument
//   - --context, -c: Context for decisions/learnings (required)
//   - --rationale, -r: Rationale for decisions (required for decisions)
//   - --consequence: Consequence for decisions (required for decisions)
//   - --lesson, -l: Lesson for learnings (required for learnings)
//   - --application, -a: Application for learnings (required for learnings)
//
// Returns:
//   - *cobra.Command: Configured add command with flags registered
func Cmd() *cobra.Command {
	var (
		priority    string
		section     string
		fromFile    string
		context     string
		rationale   string
		consequence string
		lesson      string
		application string
	)

	short, long := desc.Command(cmd.DescKeyAdd)

	c := &cobra.Command{
		Use:       cmd.UseAdd,
		Short:     short,
		Long:      long,
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{entry.Task, entry.Decision, entry.Learning, entry.Convention},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, core.Config{
				Priority:    priority,
				Section:     section,
				FromFile:    fromFile,
				Context:     context,
				Rationale:   rationale,
				Consequence: consequence,
				Lesson:      lesson,
				Application: application,
			})
		},
	}

	c.Flags().StringVarP(
		&priority,
		cflag.Priority, cflag.ShortPriority, "",
		desc.Flag(flag.DescKeyAddPriority),
	)
	_ = c.RegisterFlagCompletionFunc(
		cflag.Priority, func(_ *cobra.Command, _ []string, _ string) (
			[]string, cobra.ShellCompDirective,
		) {
			return entry.Priorities, cobra.ShellCompDirectiveNoFileComp
		})
	c.Flags().StringVarP(
		&section,
		cflag.Section, cflag.ShortSection, "",
		desc.Flag(flag.DescKeyAddSection),
	)
	c.Flags().StringVarP(
		&fromFile,
		cflag.File, cflag.ShortFile, "",
		desc.Flag(flag.DescKeyAddFile),
	)
	flagbind.StringFlagP(c, &context, cflag.Context, cflag.ShortContext, flag.DescKeyAddContext)
	flagbind.StringFlagP(c, &rationale, cflag.Rationale, cflag.ShortRationale, flag.DescKeyAddRationale)
	flagbind.StringFlag(c, &consequence, cflag.Consequence, flag.DescKeyAddConsequence)
	flagbind.StringFlagP(c, &lesson, cflag.Lesson, cflag.ShortLesson, flag.DescKeyAddLesson)
	flagbind.StringFlagP(c, &application, cflag.Application, cflag.ShortApplication, flag.DescKeyAddApplication)

	return c
}
