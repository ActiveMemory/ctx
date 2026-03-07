//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd returns the remind dismiss subcommand.
//
// Returns:
//   - *cobra.Command: Configured dismiss subcommand
func Cmd() *cobra.Command {
	var allFlag bool

	cmd := &cobra.Command{
		Use:     "dismiss [ID]",
		Aliases: []string{"rm"},
		Short:   "Dismiss one or all reminders",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return runDismissAll(cmd)
			}
			if len(args) == 0 {
				return fmt.Errorf("provide a reminder ID or use --all")
			}
			return runDismiss(cmd, args[0])
		},
	}

	cmd.Flags().BoolVar(&allFlag, "all", false, "Dismiss all reminders")

	return cmd
}
