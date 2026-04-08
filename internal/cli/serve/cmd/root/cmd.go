//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/serve/core/shared"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the serve command.
//
// Serves a static site (default) or starts the shared
// context hub when --shared is passed. Use --daemon to
// run in the background, --stop to kill a running daemon.
// Use --peers for cluster mode with Raft leader election.
//
// Returns:
//   - *cobra.Command: The serve command
func Cmd() *cobra.Command {
	var (
		isShared bool
		isDaemon bool
		isStop   bool
		port     int
		dataDir  string
		peersStr string
	)

	short, long := desc.Command(cmd.DescKeyServe)

	c := &cobra.Command{
		Use:     cmd.UseServe,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyServe),
		Args:    cobra.MaximumNArgs(1),
		RunE: func(
			cobraCmd *cobra.Command, args []string,
		) error {
			if isStop {
				return shared.Stop(
					cobraCmd, dataDir,
				)
			}

			peers := shared.ParsePeers(peersStr)

			if isShared && isDaemon {
				return shared.RunDaemon(
					cobraCmd, port, dataDir,
				)
			}
			if isShared {
				return shared.Run(
					cobraCmd, port, dataDir, peers,
				)
			}
			return Run(args)
		},
	}

	flagbind.BoolFlag(
		c, &isShared,
		cFlag.Shared, flag.DescKeyServeShared,
	)
	flagbind.IntFlag(
		c, &port,
		cFlag.Port, shared.DefaultPort(),
		flag.DescKeyServePort,
	)
	flagbind.StringFlag(
		c, &dataDir,
		cFlag.DataDir, flag.DescKeyServeDataDir,
	)
	flagbind.BoolFlag(
		c, &isDaemon,
		cFlag.Daemon, flag.DescKeyServeDaemon,
	)
	flagbind.BoolFlag(
		c, &isStop,
		cFlag.Stop, flag.DescKeyServeStop,
	)
	flagbind.StringFlag(
		c, &peersStr,
		cFlag.Peers, flag.DescKeyServePeers,
	)

	return c
}
