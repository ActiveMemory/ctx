//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/steering/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/steering/cmd/initcmd"
	"github.com/ActiveMemory/ctx/internal/cli/steering/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/steering/cmd/preview"
	"github.com/ActiveMemory/ctx/internal/cli/steering/cmd/synccmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx steering" parent command.
//
// Returns:
//   - *cobra.Command: Configured steering command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeySteering, cmd.UseSteering,
		add.Cmd(),
		list.Cmd(),
		preview.Cmd(),
		initcmd.Cmd(),
		synccmd.Cmd(),
	)
}
