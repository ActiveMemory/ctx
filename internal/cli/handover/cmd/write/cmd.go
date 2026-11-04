//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the `ctx handover write <title>` command.
//
// Returns:
//   - *cobra.Command: configured cobra command with required
//     --summary, --facts, --verified, --dead-ends, and --next
//     flags and placeholder rejection.
func Cmd() *cobra.Command {
	var (
		summary       string
		facts         string
		verified      string
		deadEnds      string
		next          string
		highlights    string
		openQuestions string
		commit        string
		noFold        bool
	)

	short, long := desc.Command(cmd.DescKeyHandoverWrite)
	c := &cobra.Command{
		Use:   cmd.UseHandoverWrite,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0], summary, facts,
				verified, deadEnds, next, highlights,
				openQuestions, commit, noFold)
		},
	}

	flagbind.StringFlag(
		c, &summary, cFlag.Summary, flag.DescKeyHandoverSummary,
	)
	flagbind.StringFlag(
		c, &facts, cFlag.Facts, flag.DescKeyHandoverFacts,
	)
	flagbind.StringFlag(
		c, &verified, cFlag.Verified,
		flag.DescKeyHandoverVerified,
	)
	flagbind.StringFlag(
		c, &deadEnds, cFlag.DeadEnds,
		flag.DescKeyHandoverDeadEnds,
	)
	flagbind.StringFlag(
		c, &next, cFlag.Next, flag.DescKeyHandoverNext,
	)
	flagbind.StringFlag(
		c, &highlights, cFlag.Highlights,
		flag.DescKeyHandoverHighlights,
	)
	flagbind.StringFlag(
		c, &openQuestions, cFlag.OpenQuestions,
		flag.DescKeyHandoverOpenQuestions,
	)
	flagbind.StringFlag(
		c, &commit, cFlag.Commit, flag.DescKeyHandoverCommit,
	)
	flagbind.BoolFlag(
		c, &noFold, cFlag.NoFold, flag.DescKeyHandoverNoFold,
	)

	// Acceptable discard: MarkFlagRequired only errors on an
	// unregistered flag name; all five flags are bound above.
	_ = c.MarkFlagRequired(cFlag.Summary)
	_ = c.MarkFlagRequired(cFlag.Facts)
	_ = c.MarkFlagRequired(cFlag.Verified)
	_ = c.MarkFlagRequired(cFlag.DeadEnds)
	_ = c.MarkFlagRequired(cFlag.Next)

	return c
}
