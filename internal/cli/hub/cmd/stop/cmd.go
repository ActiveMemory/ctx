//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stop

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

// Cmd returns the hub stop subcommand.
//
// Sends SIGTERM to a daemonized ctx Hub server using the PID
// file under the hub data directory, then removes the PID file.
//
// Returns:
//   - *cobra.Command: The stop subcommand
func Cmd() *cobra.Command {
	var dataDir string

	short, long := desc.Command(cmd.DescKeyHubStop)

	c := &cobra.Command{
		Use:     cmd.UseHubStop,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubStop),
		Args:    cobra.NoArgs,
		// Hub stores at ~/.ctx/hub-data/, not .context/.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			return server.Stop(cobraCmd, dataDir)
		},
	}

	flagbind.StringFlag(
		c, &dataDir,
		cFlag.DataDir, flag.DescKeyHubStopDataDir,
	)

	return c
}
