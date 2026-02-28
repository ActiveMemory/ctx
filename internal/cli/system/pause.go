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

// pauseCmd returns the "ctx system pause" plumbing command.
//
// Creates a session-scoped pause marker that causes all nudge hooks to
// no-op. Security hooks (block-non-path-ctx, block-dangerous-commands)
// and housekeeping hooks (cleanup-tmp) are unaffected.
func pauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause context hooks for this session",
		Long: `Creates a session-scoped pause marker. While paused, all nudge
and reminder hooks no-op. Security and housekeeping hooks still fire.

The session ID is read from stdin JSON (same as hooks) or --session-id flag.`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPause(cmd, os.Stdin)
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}

func runPause(cmd *cobra.Command, stdin *os.File) error {
	sessionID, _ := cmd.Flags().GetString("session-id")
	if sessionID == "" {
		input := readInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = sessionUnknown
	}

	path := pauseMarkerPath(sessionID)
	writeCounter(path, 0)
	cmd.Printf("Context hooks paused for session %s\n", sessionID)
	return nil
}
