//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package build

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/add/core/run"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd builds a cobra add subcommand for the given noun.
//
// Parameters:
//   - noun: One of entry.{Task,Decision,Learning,Convention}.
//     Prepended to args before invoking run.Run.
//   - descKey: Description key for the embedded asset lookup
//     (e.g., "task.add", "decision.add").
//   - useStr: Cobra Use string (typically "add [content]").
//
// Returns:
//   - *cobra.Command: Configured add subcommand with all flags
//     registered.
func Cmd(noun, descKey, useStr string) *cobra.Command {
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

	short, long := desc.Command(descKey)

	c := &cobra.Command{
		Use:     useStr,
		Short:   short,
		Long:    long,
		Example: desc.Example(descKey),
		RunE: func(cmd *cobra.Command, args []string) error {
			withNoun := append([]string{noun}, args...)
			return run.Run(cmd, withNoun, entity.AddConfig{
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

	flagbind.BindStringFlagsP(c,
		[]*string{
			&priority, &section, &fromFile, &context,
			&rationale, &lesson, &application,
		},
		[]string{
			cFlag.Priority, cFlag.Section, cFlag.File, cFlag.Context,
			cFlag.Rationale, cFlag.Lesson, cFlag.Application,
		},
		[]string{
			cFlag.ShortPriority, cFlag.ShortSection,
			cFlag.ShortFile, cFlag.ShortContext,
			cFlag.ShortRationale, cFlag.ShortLesson,
			cFlag.ShortApplication,
		},
		[]string{
			flag.DescKeyAddPriority, flag.DescKeyAddSection,
			flag.DescKeyAddFile, flag.DescKeyAddContext,
			flag.DescKeyAddRationale, flag.DescKeyAddLesson,
			flag.DescKeyAddApplication,
		},
	)
	flagbind.BindStringFlags(c,
		[]*string{
			&consequence, &sessionID, &branch, &commit,
		},
		[]string{
			cFlag.Consequence, cFlag.SessionID,
			cFlag.Branch, cFlag.Commit,
		},
		[]string{
			flag.DescKeyAddConsequence,
			flag.DescKeyAddSessionID,
			flag.DescKeyAddBranch, flag.DescKeyAddCommit,
		},
	)
	flagbind.BoolFlag(
		c, &share,
		cFlag.Share, flag.DescKeyAddShare,
	)

	_ = c.RegisterFlagCompletionFunc(
		cFlag.Priority, func(
			_ *cobra.Command, _ []string, _ string,
		) ([]string, cobra.ShellCompDirective) {
			return cfgEntry.Priorities,
				cobra.ShellCompDirectiveNoFileComp
		})

	return c
}
