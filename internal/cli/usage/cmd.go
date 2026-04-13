//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package usage

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	cfgStats "github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx usage" top-level command.
//
// Returns:
//   - *cobra.Command: Configured usage command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyUsage)

	c := &cobra.Command{
		Use:     cmd.UseUsage,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyUsage),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	flagbind.BoolFlagShort(c,
		cFlag.Follow, cFlag.ShortFollow,
		flag.DescKeyUsageFollow,
	)
	flagbind.StringFlagShort(c,
		cFlag.Session, cFlag.ShortSessionID,
		flag.DescKeyUsageSession,
	)
	flagbind.LastJSON(c, cfgStats.DefaultLast,
		flag.DescKeyUsageLast,
		flag.DescKeyUsageJson,
	)

	return c
}
