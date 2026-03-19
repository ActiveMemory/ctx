//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the top-level "ctx pause" command.
//
// Returns:
//   - *cobra.Command: Configured pause command
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyPause)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyPause,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString("session-id")
			return Run(cmd, sessionID)
		},
	}
	cmd.Flags().String("session-id", "",
		desc.FlagDesc(flag.DescKeyPauseSessionId),
	)
	return cmd
}
