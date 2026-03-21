//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mv

import (
	"strconv"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the pad mv subcommand.
//
// Returns:
//   - *cobra.Command: Configured mv subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyPadMv)
	return &cobra.Command{
		Use:   cmd.UsePadMv,
		Short: short,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			m, cErr := strconv.Atoi(args[1])
			if cErr != nil {
				return cErr
			}
			return Run(cmd, n, m)
		},
	}
}
