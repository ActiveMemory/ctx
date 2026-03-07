//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"github.com/spf13/cobra"

	resumeroot "github.com/ActiveMemory/ctx/internal/cli/resume/cmd/root"
)

// Cmd returns the top-level "ctx resume" command.
//
// Returns:
//   - *cobra.Command: Configured resume command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume context hooks for this session",
		Long: `Resume context hooks after a pause. Silent no-op if not paused.

The session ID is read from stdin JSON (same as hooks) or --session-id flag.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString("session-id")
			return resumeroot.Run(cmd, sessionID)
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}
