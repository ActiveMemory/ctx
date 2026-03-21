//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/edit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/reset"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/show"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system message" subcommand.
//
// Returns:
//   - *cobra.Command: Configured message subcommand with sub-subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemMessage)

	c := &cobra.Command{
		Use:   cmd.UseSystemMessage,
		Short: short,
		Long:  long,
	}

	c.AddCommand(
		list.Cmd(),
		show.Cmd(),
		edit.Cmd(),
		reset.Cmd(),
	)

	return c
}
