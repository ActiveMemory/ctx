//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
)

// Cmd returns the recall list subcommand.
//
// Returns:
//   - *cobra.Command: Command for listing parsed sessions
func Cmd() *cobra.Command {
	var (
		limit       int
		project     string
		tool        string
		since       string
		until       string
		allProjects bool
	)

	short, long := desc.Command(cmd.DescKeyRecallList)

	c := &cobra.Command{
		Use:   cmd.UseRecallList,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, limit, project, tool, since, until, allProjects)
		},
	}

	c.Flags().IntVarP(
		&limit, cFlag.Limit,
		cFlag.ShortMaxIterations, journal.DefaultRecallListLimit,
		desc.Flag(flag.DescKeyRecallListLimit),
	)
	c.Flags().StringVarP(&project, "project", "p", "",
		desc.Flag(flag.DescKeyRecallListProject),
	)
	c.Flags().StringVarP(&tool, cFlag.Tool, cFlag.ShortTool, "",
		desc.Flag(flag.DescKeyRecallListTool),
	)
	c.Flags().StringVar(&since, cFlag.Since, "",
		desc.Flag(flag.DescKeyRecallListSince),
	)
	c.Flags().StringVar(&until, cFlag.Until, "",
		desc.Flag(flag.DescKeyRecallListUntil),
	)
	c.Flags().BoolVar(&allProjects, cFlag.AllProjects, false,
		desc.Flag(flag.DescKeyRecallListAllProjects),
	)

	return c
}
