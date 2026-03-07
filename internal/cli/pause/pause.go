//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	pauseroot "github.com/ActiveMemory/ctx/internal/cli/pause/cmd/root"
)

// Cmd returns the top-level "ctx pause" command.
//
// Returns:
//   - *cobra.Command: Configured pause command
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("pause")
	cmd := &cobra.Command{
		Use:   "pause",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString("session-id")
			return pauseroot.Run(cmd, sessionID)
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}
