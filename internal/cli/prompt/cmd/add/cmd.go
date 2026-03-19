//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the prompt add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var fromStdin bool

	short, long := desc.CommandDesc(cmd.DescKeyPromptAdd)

	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0], fromStdin)
		},
	}

	cmd.Flags().BoolVar(&fromStdin,
		"stdin", false, desc.FlagDesc(flag.DescKeyPromptAddStdin),
	)

	return cmd
}
