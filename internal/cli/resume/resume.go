//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resume provides the top-level "ctx resume" command.
//
// Removes the session-scoped pause marker, restoring normal hook behavior.
// Silent no-op if the session is not currently paused.
package resume

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system"
)

// Cmd returns the top-level "ctx resume" command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume context hooks for this session",
		Long: `Resume context hooks after a pause. Silent no-op if not paused.

The session ID is read from stdin JSON (same as hooks) or --session-id flag.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString("session-id")
			if sessionID == "" {
				sessionID = system.ReadSessionID(os.Stdin)
			}
			system.Resume(sessionID)
			cmd.Printf("Context hooks resumed for session %s\n", sessionID)
			return nil
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}
