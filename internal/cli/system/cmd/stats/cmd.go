//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx system stats" subcommand.
//
// Returns:
//   - *cobra.Command: Configured stats subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemStats)

	c := &cobra.Command{
		Use:   cmd.UseSystemStats,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	c.Flags().BoolP(cFlag.Follow, cFlag.ShortFollow, false,
		desc.Flag(flag.DescKeySystemStatsFollow),
	)
	c.Flags().StringP(cFlag.Session, cFlag.ShortSessionID, "",
		desc.Flag(flag.DescKeySystemStatsSession),
	)
	flagbind.LastJSON(c, stats.DefaultLast,
		flag.DescKeySystemStatsLast, flag.DescKeySystemStatsJson,
	)

	return c
}
