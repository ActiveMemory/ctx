//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreStatus "github.com/ActiveMemory/ctx/internal/cli/connection/core/status"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the connect status subcommand.
//
// Returns:
//   - *cobra.Command: The status subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyConnectionStatus)

	return &cobra.Command{
		Use:     cmd.UseConnectionStatus,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyConnectionStatus),
		Args:    cobra.NoArgs,
		RunE:    coreStatus.Run,
	}
}
