//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the changes command.
//
// Returns:
//   - *cobra.Command: Configured changes command with flags registered
func Cmd() *cobra.Command {
	var since string

	short, long := assets.CommandDesc("changes")

	cmd := &cobra.Command{
		Use:   "changes",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, since)
		},
	}

	cmd.Flags().StringVar(&since, "since", "", assets.FlagDesc("changes.since"))

	return cmd
}
