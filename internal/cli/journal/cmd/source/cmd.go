//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreSrc "github.com/ActiveMemory/ctx/internal/cli/journal/core/source"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the journal source subcommand.
//
// Combines session listing and inspection into a single entry
// point. Default behavior (no flags) lists available sessions.
// Use --show to inspect a specific session by slug or ID.
//
// Returns:
//   - *cobra.Command: Command for listing and inspecting
//     session sources
func Cmd() *cobra.Command {
	var (
		showID      string
		latest      bool
		full        bool
		limit       int
		project     string
		tool        string
		since       string
		until       string
		allProjects bool
	)

	short, long := desc.Command(cmd.DescKeyJournalSource)

	c := &cobra.Command{
		Use:     cmd.UseJournalSource,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalSource),
		RunE: func(
			cmd *cobra.Command, args []string,
		) error {
			return Run(cmd, args, coreSrc.Opts{
				ShowID:      showID,
				Latest:      latest,
				Full:        full,
				Limit:       limit,
				Project:     project,
				Tool:        tool,
				Since:       since,
				Until:       until,
				AllProjects: allProjects,
			})
		},
	}

	flagbind.StringFlagP(
		c, &showID,
		cFlag.Show, cFlag.ShortShow,
		flag.DescKeyJournalSourceShow,
	)
	flagbind.BoolFlag(
		c, &latest,
		cFlag.Latest, flag.DescKeyJournalSourceLatest,
	)
	flagbind.BoolFlag(
		c, &full,
		cFlag.Full, flag.DescKeyJournalSourceFull,
	)
	flagbind.IntFlagP(
		c, &limit,
		cFlag.Limit, cFlag.ShortMaxIterations,
		journal.DefaultRecallListLimit,
		flag.DescKeyJournalSourceLimit,
	)
	flagbind.StringFlagP(
		c, &project,
		cFlag.Project, cFlag.ShortProject,
		flag.DescKeyJournalSourceProject,
	)
	flagbind.StringFlagP(
		c, &tool,
		cFlag.Tool, cFlag.ShortTool,
		flag.DescKeyJournalSourceTool,
	)
	flagbind.StringFlag(
		c, &since,
		cFlag.Since, flag.DescKeyJournalSourceSince,
	)
	flagbind.StringFlag(
		c, &until,
		cFlag.Until, flag.DescKeyJournalSourceUntil,
	)
	flagbind.BoolFlag(
		c, &allProjects,
		cFlag.AllProjects,
		flag.DescKeyJournalSourceAllProjects,
	)

	return c
}
