//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"os"

	session2 "github.com/ActiveMemory/ctx/internal/write/session"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/session"
)

// Run executes the resume logic.
//
// Reads a session ID from the --session-id flag or stdin JSON, then
// removes the pause marker file so hooks fire normally again.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil
func Run(cmd *cobra.Command, stdin *os.File) error {
	sessionID, _ := cmd.Flags().GetString("session-id")
	if sessionID == "" {
		input := core.ReadInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = session.IDUnknown
	}

	path := core.PauseMarkerPath(sessionID)
	_ = os.Remove(path)
	session2.SessionResumed(cmd, sessionID)
	return nil
}
