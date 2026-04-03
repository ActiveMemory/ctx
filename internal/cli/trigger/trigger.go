//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/trigger/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/trigger/cmd/disable"
	"github.com/ActiveMemory/ctx/internal/cli/trigger/cmd/enable"
	"github.com/ActiveMemory/ctx/internal/cli/trigger/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/trigger/cmd/test"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx trigger" parent command.
//
// Returns:
//   - *cobra.Command: Configured trigger command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyTrigger, cmd.UseTrigger,
		add.Cmd(),
		list.Cmd(),
		test.Cmd(),
		enable.Cmd(),
		disable.Cmd(),
	)
}
