//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	embedCmd "github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx disclosure inspect" command.
//
// It takes one positional argument (a canonical knowledge file) and
// reports its staged entries and current themes. Read-only; --json
// switches to machine-readable output.
//
// Returns:
//   - *cobra.Command: configured inspect command with --json registered
func Cmd() *cobra.Command {
	var jsonOutput bool

	short, long := desc.Command(embedCmd.DescKeyDisclosureInspect)
	c := &cobra.Command{
		Use:   embedCmd.UseDisclosureInspect,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], jsonOutput)
		},
	}

	flagbind.BoolFlag(c, &jsonOutput, cFlag.JSON, flag.DescKeyIndexJSON)
	return c
}
