//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package topic

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/topic/cmd/newcmd"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx kb topic` parent command with the
// `new` subcommand registered.
//
// Returns:
//   - *cobra.Command: topic parent with `new` registered.
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyKBTopic, cmd.UseKBTopic,
		newcmd.Cmd(),
	)
}
