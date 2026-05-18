//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	writeCmd "github.com/ActiveMemory/ctx/internal/cli/handover/cmd/write"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx handover` parent command.
//
// Returns:
//   - *cobra.Command: parent with `write` subcommand registered.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyHandover)
	c := &cobra.Command{
		Use:   cmd.UseHandover,
		Short: short,
		Long:  long,
	}
	c.AddCommand(writeCmd.Cmd())
	return c
}
