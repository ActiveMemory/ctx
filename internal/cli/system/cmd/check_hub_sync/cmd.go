//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_hub_sync

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the check-hub-sync hook command.
//
// Returns:
//   - *cobra.Command: The hook command
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:    "check-hub-sync",
		Short:  "Auto-sync shared hub entries",
		Hidden: true,
		RunE: func(
			cmd *cobra.Command, _ []string,
		) error {
			return Run(cmd, os.Stdin)
		},
	}
}
