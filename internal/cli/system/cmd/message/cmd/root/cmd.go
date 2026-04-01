//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/edit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/reset"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/show"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system message" subcommand.
//
// Returns:
//   - *cobra.Command: Configured message command
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeySystemMessage, cmd.UseSystemMessage,
		list.Cmd(),
		show.Cmd(),
		edit.Cmd(),
		reset.Cmd(),
	)
}
