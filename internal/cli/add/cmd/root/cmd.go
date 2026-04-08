//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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
		sessionID   string
		branch      string
		commit      string
		context     string
		rationale   string
		consequence string
		lesson      string
		application string
		share       bool
	)

	short, long := desc.Command(cmd.DescKeyAdd)

	c := &cobra.Command{
		Use:     cmd.UseAdd,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyAdd),
		Args:    cobra.MinimumNArgs(1),
		ValidArgs: []string{
			entry.Task,
			entry.Decision,
			entry.Learning,
			entry.Convention,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, entity.AddConfig{
				Priority:    priority,
				Section:     section,
				FromFile:    fromFile,
				SessionID:   sessionID,
				Branch:      branch,
				Commit:      commit,
				Context:     context,
				Rationale:   rationale,
				Consequence: consequence,
				Lesson:      lesson,
				Application: application,
				Share:       share,
			})
		},
	}

	flagbind.StringFlagP(
		c, &priority,
		cFlag.Priority, cFlag.ShortPriority,
		flag.DescKeyAddPriority,
	)
	_ = c.RegisterFlagCompletionFunc(
		cFlag.Priority, func(_ *cobra.Command, _ []string, _ string) (
			[]string, cobra.ShellCompDirective,
		) {
			return entry.Priorities, cobra.ShellCompDirectiveNoFileComp
		})
	flagbind.StringFlagP(
		c, &section,
		cFlag.Section, cFlag.ShortSection,
		flag.DescKeyAddSection,
	)
	flagbind.StringFlagP(
		c, &fromFile,
		cFlag.File, cFlag.ShortFile,
		flag.DescKeyAddFile,
	)
	flagbind.StringFlagP(
		c, &context,
		cFlag.Context, cFlag.ShortContext, flag.DescKeyAddContext,
	)
	flagbind.StringFlagP(
		c, &rationale,
		cFlag.Rationale, cFlag.ShortRationale, flag.DescKeyAddRationale,
	)
	flagbind.StringFlag(
		c, &consequence,
		cFlag.Consequence, flag.DescKeyAddConsequence,
	)
	flagbind.StringFlagP(
		c, &lesson,
		cFlag.Lesson, cFlag.ShortLesson, flag.DescKeyAddLesson,
	)
	flagbind.StringFlagP(
		c, &application,
		cFlag.Application, cFlag.ShortApplication, flag.DescKeyAddApplication,
	)
	flagbind.StringFlag(
		c, &sessionID,
		cFlag.SessionID, flag.DescKeyAddSessionID,
	)
	flagbind.StringFlag(
		c, &branch,
		cFlag.Branch, flag.DescKeyAddBranch,
	)
	flagbind.StringFlag(
		c, &commit,
		cFlag.Commit, flag.DescKeyAddCommit,
	)
	flagbind.BoolFlag(
		c, &share,
		cFlag.Share, flag.DescKeyAddShare,
	)

	return c
}
