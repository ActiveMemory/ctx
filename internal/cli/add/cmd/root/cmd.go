//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core"
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

	short, long := desc.CommandDesc(cmd.DescKeyAdd)

	cmd := &cobra.Command{
		Use:       cmd.DescKeyAdd + " <type> [content]",
		Short:     short,
		Long:      long,
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"task", "decision", "learning", "convention"},
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

	cmd.Flags().StringVarP(
		&priority,
		"priority", "p", "",
		desc.FlagDesc(flag.DescKeyAddPriority),
	)
	_ = cmd.RegisterFlagCompletionFunc(
		"priority", func(_ *cobra.Command, _ []string, _ string) (
			[]string, cobra.ShellCompDirective,
		) {
			return []string{"high", "medium", "low"}, cobra.ShellCompDirectiveNoFileComp
		})
	cmd.Flags().StringVarP(
		&section,
		"section", "s", "",
		desc.FlagDesc(flag.DescKeyAddSection),
	)
	cmd.Flags().StringVarP(
		&fromFile,
		"file", "f", "",
		desc.FlagDesc(flag.DescKeyAddFile),
	)
	cmd.Flags().StringVarP(
		&context,
		"context", "c", "",
		desc.FlagDesc(flag.DescKeyAddContext),
	)
	cmd.Flags().StringVarP(
		&rationale,
		"rationale", "r", "",
		desc.FlagDesc(flag.DescKeyAddRationale),
	)
	cmd.Flags().StringVar(
		&consequence,
		"consequence", "",
		desc.FlagDesc(flag.DescKeyAddConsequence),
	)
	cmd.Flags().StringVarP(
		&lesson,
		"lesson", "l", "",
		desc.FlagDesc(flag.DescKeyAddLesson),
	)
	cmd.Flags().StringVarP(
		&application,
		"application", "a", "",
		desc.FlagDesc(flag.DescKeyAddApplication),
	)

	return cmd
}
