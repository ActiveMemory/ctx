//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package listen

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreListen "github.com/ActiveMemory/ctx/internal/cli/connection/core/listen"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the connect listen subcommand.
//
// Returns:
//   - *cobra.Command: The listen subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyConnectionListen)

	return &cobra.Command{
		Use:     cmd.UseConnectionListen,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyConnectionListen),
		Args:    cobra.NoArgs,
		RunE:    coreListen.Run,
	}
}
