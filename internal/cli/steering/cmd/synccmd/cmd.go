//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package synccmd

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/resolve"
	coreSync "github.com/ActiveMemory/ctx/internal/cli/steering/core/sync"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
)

// Cmd returns the "ctx steering sync" subcommand.
//
// Returns:
//   - *cobra.Command: Configured sync subcommand
func Cmd() *cobra.Command {
	var syncAll bool

	short, long := desc.Command(cmd.DescKeySteeringSync)

	c := &cobra.Command{
		Use:   cmd.UseSteeringSync,
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c, syncAll)
		},
	}

	flagbind.BoolFlag(
		c, &syncAll, flag.All,
		embedFlag.DescKeySteeringSyncAll,
	)

	return c
}

// Run syncs steering files to tool-native formats.
//
// Parameters:
//   - c: The cobra command for output and flag access
//   - syncAll: Whether to sync to all supported tools
func Run(c *cobra.Command, syncAll bool) error {
	steeringDir := rc.SteeringDir()
	projectRoot := token.Dot

	if syncAll {
		report, syncErr := steering.SyncAll(
			steeringDir, projectRoot,
		)
		if syncErr != nil {
			return syncErr
		}
		coreSync.PrintReport(c, report)
		return nil
	}

	// Resolve tool from --tool flag or .ctxrc.
	tool, resolveErr := resolve.Tool(c)
	if resolveErr != nil {
		return errSteering.NoTool()
	}

	report, syncErr := steering.SyncTool(
		steeringDir, projectRoot, tool,
	)
	if syncErr != nil {
		return syncErr
	}

	coreSync.PrintReport(c, report)
	return nil
}
