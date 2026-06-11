//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/dream/cmd/accept"
	"github.com/ActiveMemory/ctx/internal/cli/dream/cmd/amend"
	"github.com/ActiveMemory/ctx/internal/cli/dream/cmd/reject"
	"github.com/ActiveMemory/ctx/internal/cli/dream/cmd/review"
	dreamPass "github.com/ActiveMemory/ctx/internal/cli/dream/core/pass"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the dream command with its subcommands.
//
// Invoked with no subcommand, it runs one dream pass; flag values fall
// back to the dream rc section, then to the config defaults.
//
// Returns:
//   - *cobra.Command: the dream command with review/accept/reject/amend
//     subcommands
func Cmd() *cobra.Command {
	var (
		mode   string
		maxN   int
		budget int
		force  bool
	)

	short, long := desc.Command(cmd.DescKeyDream)
	c := &cobra.Command{
		Use:     cmd.UseDream,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyDream),
		Args:    cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			return dreamPass.Run(cobraCmd, dreamPass.Opts{
				Mode:   mode,
				Max:    maxN,
				Budget: budget,
				Force:  force,
			})
		},
	}

	flagbind.StringFlagDefault(c, &mode,
		cFlag.Mode, rc.DreamMode(), flag.DescKeyDreamMode,
	)
	flagbind.IntFlag(c, &maxN,
		cFlag.Max, rc.DreamMax(), flag.DescKeyDreamMax,
	)
	flagbind.IntFlag(c, &budget,
		cFlag.Budget, rc.DreamBudget(), flag.DescKeyDreamBudget,
	)
	flagbind.BoolFlag(c, &force, cFlag.Force, flag.DescKeyDreamForce)

	c.AddCommand(review.Cmd())
	c.AddCommand(accept.Cmd())
	c.AddCommand(reject.Cmd())
	c.AddCommand(amend.Cmd())
	return c
}
