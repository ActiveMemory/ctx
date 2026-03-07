//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"github.com/spf13/cobra"

	pauseroot "github.com/ActiveMemory/ctx/internal/cli/pause/cmd/root"
)

// Cmd returns the top-level "ctx pause" command.
//
// Returns:
//   - *cobra.Command: Configured pause command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause context hooks for this session",
		Long: `Pause all context nudge and reminder hooks for the current session.
Security hooks (dangerous command blocking) and housekeeping hooks still fire.

The session ID is read from stdin JSON (same as hooks) or --session-id flag.
Resume with: ctx resume`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString("session-id")
			return pauseroot.Run(cmd, sessionID)
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}
