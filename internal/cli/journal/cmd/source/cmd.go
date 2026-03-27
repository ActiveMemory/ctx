//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
)

// Cmd returns the journal source subcommand.
//
// Combines session listing and inspection into a single entry point.
// Default behavior (no flags) lists available sessions. Use --show to
// inspect a specific session by slug or ID.
//
// Returns:
//   - *cobra.Command: Command for listing and inspecting session sources
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
		Use:   cmd.UseJournalSource,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, Opts{
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

	c.Flags().StringVarP(&showID, cFlag.Show, cFlag.ShortShow, "",
		desc.Flag(flag.DescKeyJournalSourceShow),
	)
	c.Flags().BoolVar(&latest, cFlag.Latest, false,
		desc.Flag(flag.DescKeyJournalSourceLatest),
	)
	c.Flags().BoolVar(&full, cFlag.Full, false,
		desc.Flag(flag.DescKeyJournalSourceFull),
	)
	c.Flags().IntVarP(
		&limit, cFlag.Limit,
		cFlag.ShortMaxIterations, journal.DefaultRecallListLimit,
		desc.Flag(flag.DescKeyJournalSourceLimit),
	)
	c.Flags().StringVarP(&project, "project", "p", "",
		desc.Flag(flag.DescKeyJournalSourceProject),
	)
	c.Flags().StringVarP(&tool, cFlag.Tool, cFlag.ShortTool, "",
		desc.Flag(flag.DescKeyJournalSourceTool),
	)
	c.Flags().StringVar(&since, cFlag.Since, "",
		desc.Flag(flag.DescKeyJournalSourceSince),
	)
	c.Flags().StringVar(&until, cFlag.Until, "",
		desc.Flag(flag.DescKeyJournalSourceUntil),
	)
	c.Flags().BoolVar(&allProjects, cFlag.AllProjects, false,
		desc.Flag(flag.DescKeyJournalSourceAllProjects),
	)

	return c
}
