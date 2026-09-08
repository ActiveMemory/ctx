//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/hub/core/server"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the hub start subcommand.
//
// Starts the ctx Hub gRPC server either in the foreground or
// as a detached daemon. When --raft-bind is set, runs a Raft
// node for leader election, bootstrapped with the nodes named
// by --peers, or waiting to be added by a leader with --join.
//
// Returns:
//   - *cobra.Command: The start subcommand
func Cmd() *cobra.Command {
	var (
		isDaemon bool
		port     int
		dataDir  string
		peersStr string
		raftBind string
		join     bool
	)

	short, long := desc.Command(cmd.DescKeyHubStart)

	c := &cobra.Command{
		Use:     cmd.UseHubStart,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubStart),
		Args:    cobra.NoArgs,
		// Hub stores at ~/.ctx/hub-data/, never reads .context/.
		// Exempt from the require-context-dir gate so AWS/EKS hub
		// users hit no-broken-windows on first contact.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			if isDaemon {
				return server.RunDaemon(
					cobraCmd, port, dataDir,
					raftBind, peersStr, join,
				)
			}
			peers := server.ParsePeers(peersStr)
			return server.Run(
				cobraCmd, port, dataDir,
				raftBind, peers, join,
			)
		},
	}

	flagbind.IntFlag(
		c, &port,
		cFlag.Port, server.DefaultPort(),
		flag.DescKeyHubStartPort,
	)
	flagbind.StringFlag(
		c, &dataDir,
		cFlag.DataDir, flag.DescKeyHubStartDataDir,
	)
	flagbind.BoolFlag(
		c, &isDaemon,
		cFlag.Daemon, flag.DescKeyHubStartDaemon,
	)
	flagbind.StringFlag(
		c, &peersStr,
		cFlag.Peers, flag.DescKeyHubStartPeers,
	)
	flagbind.StringFlag(
		c, &raftBind,
		cFlag.RaftBind, flag.DescKeyHubStartRaftBind,
	)
	flagbind.BoolFlag(
		c, &join,
		cFlag.Join, flag.DescKeyHubStartJoin,
	)

	return c
}
