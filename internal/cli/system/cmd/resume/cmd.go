//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Cmd returns the "ctx system resume" plumbing command.
//
// Returns:
//   - *cobra.Command: Configured resume subcommand
func Cmd() *cobra.Command {
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
		input := core.ReadInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}

	path := core.PauseMarkerPath(sessionID)
	_ = os.Remove(path)
	cmd.Println(fmt.Sprintf("Context hooks resumed for session %s", sessionID))
	return nil
}
