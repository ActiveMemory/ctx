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
	cfgStats "github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx stats" top-level command.
//
// Returns:
//   - *cobra.Command: Configured stats command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyStats)

	c := &cobra.Command{
		Use:     cmd.UseStats,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyStats),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	flagbind.BoolFlagShort(c,
		cFlag.Follow, cFlag.ShortFollow,
		flag.DescKeyStatsFollow,
	)
	flagbind.StringFlagShort(c,
		cFlag.Session, cFlag.ShortSessionID,
		flag.DescKeyStatsSession,
	)
	flagbind.LastJSON(c, cfgStats.DefaultLast,
		flag.DescKeyStatsLast,
		flag.DescKeyStatsJson,
	)

	return c
}
