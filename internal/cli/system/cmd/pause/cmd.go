//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Cmd returns the "ctx system pause" plumbing command.
//
// Returns:
//   - *cobra.Command: Configured pause subcommand
func Cmd() *cobra.Command {
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
		input := core.ReadInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}

	path := core.PauseMarkerPath(sessionID)
	core.WriteCounter(path, 0)
	cmd.Println(fmt.Sprintf("Context hooks paused for session %s", sessionID))
	return nil
}
