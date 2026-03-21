//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the prompt add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var fromStdin bool

	short, long := desc.Command(cmd.DescKeyPromptAdd)

	c := &cobra.Command{
		Use:   cmd.UsePromptAdd,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0], fromStdin)
		},
	}

	c.Flags().BoolVar(&fromStdin,
		cflag.Stdin, false, desc.Flag(flag.DescKeyPromptAddStdin),
	)

	return c
}
