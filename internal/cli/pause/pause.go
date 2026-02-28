//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pause provides the top-level "ctx pause" command.
//
// Pauses all context nudge hooks for the current session. Security and
// housekeeping hooks are unaffected. Delegates to the session-scoped
// pause marker in secureTempDir.
package pause

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system"
)

// Cmd returns the top-level "ctx pause" command.
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
			if sessionID == "" {
				sessionID = system.ReadSessionID(os.Stdin)
			}
			system.Pause(sessionID)
			cmd.Printf("Context hooks paused for session %s\n", sessionID)
			return nil
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}
