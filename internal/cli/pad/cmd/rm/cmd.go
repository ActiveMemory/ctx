//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"strconv"

	"github.com/spf13/cobra"
)

// Cmd returns the pad rm subcommand.
//
// Returns:
//   - *cobra.Command: Configured rm subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm N",
		Short: "Remove an entry by number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return runRm(cmd, n)
		},
	}
}
