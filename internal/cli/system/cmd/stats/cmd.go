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
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/stats"
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

	c.Flags().BoolP(cflag.Follow, cflag.ShortFollow, false,
		desc.Flag(flag.DescKeySystemStatsFollow),
	)
	c.Flags().StringP(cflag.Session, cflag.ShortSessionID, "",
		desc.Flag(flag.DescKeySystemStatsSession),
	)
	c.Flags().IntP(cflag.Last, cflag.ShortLast, stats.DefaultLast,
		desc.Flag(flag.DescKeySystemStatsLast),
	)
	c.Flags().BoolP(cflag.JSON, cflag.ShortJSON, false,
		desc.Flag(flag.DescKeySystemStatsJson),
	)

	return c
}
