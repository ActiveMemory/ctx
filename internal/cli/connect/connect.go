//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package connect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/connect/cmd/register"
	"github.com/ActiveMemory/ctx/internal/cli/connect/cmd/subscribe"
	connectSync "github.com/ActiveMemory/ctx/internal/cli/connect/cmd/sync"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the connect command with subcommands.
//
// Returns:
//   - *cobra.Command: connect with register, subscribe, sync
func Cmd() *cobra.Command {
	return parent.Cmd(
		cmd.DescKeyConnect, cmd.UseConnect,
		register.Cmd(),
		subscribe.Cmd(),
		connectSync.Cmd(),
	)
}
