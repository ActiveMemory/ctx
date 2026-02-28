//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"

	"github.com/spf13/cobra"
)

// resumeCmd returns the "ctx system resume" plumbing command.
//
// Removes the session-scoped pause marker, restoring normal hook behavior.
// Silent no-op if the session is not paused.
func resumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume context hooks for this session",
		Long: `Removes the session-scoped pause marker. Hooks resume normal
behavior. Silent no-op if not paused.

The session ID is read from stdin JSON (same as hooks) or --session-id flag.`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResume(cmd, os.Stdin)
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}

func runResume(cmd *cobra.Command, stdin *os.File) error {
	sessionID, _ := cmd.Flags().GetString("session-id")
	if sessionID == "" {
		input := readInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = sessionUnknown
	}

	path := pauseMarkerPath(sessionID)
	_ = os.Remove(path)
	cmd.Printf("Context hooks resumed for session %s\n", sessionID)
	return nil
}
