//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the recall show subcommand.
//
// Returns:
//   - *cobra.Command: Command for showing session details
func Cmd() *cobra.Command {
	var (
		latest      bool
		full        bool
		allProjects bool
	)

	short, long := desc.Command(cmd.DescKeyRecallShow)

	c := &cobra.Command{
		Use:   cmd.UseRecallShow,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, latest, full, allProjects)
		},
	}

	c.Flags().BoolVar(&latest, cFlag.Latest, false,
		desc.Flag(flag.DescKeyRecallShowLatest),
	)
	c.Flags().BoolVar(&full, cFlag.Full, false,
		desc.Flag(flag.DescKeyRecallShowFull),
	)
	c.Flags().BoolVar(&allProjects, cFlag.AllProjects, false,
		desc.Flag(flag.DescKeyRecallShowAllProjects),
	)

	return c
}
