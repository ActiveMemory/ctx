//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/rm"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/show"
)

// Cmd returns the prompt command with subcommands.
//
// When invoked without a subcommand, it lists all prompt templates.
//
// Returns:
//   - *cobra.Command: Configured prompt command with subcommands
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("prompt")

	cmd := &cobra.Command{
		Use:   "prompt",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return list.Run(cmd)
		},
	}

	cmd.AddCommand(list.Cmd())
	cmd.AddCommand(show.Cmd())
	cmd.AddCommand(add.Cmd())
	cmd.AddCommand(rm.Cmd())

	return cmd
}
