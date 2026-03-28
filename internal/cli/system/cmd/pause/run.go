//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/session"
	writePause "github.com/ActiveMemory/ctx/internal/write/pause"
)

// Run executes the pause logic.
//
// Reads a session ID from the --session-id flag or stdin JSON, then
// creates a pause marker file so all subsequent hooks for that session
// are suppressed until resumed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil
func Run(cmd *cobra.Command, stdin *os.File) error {
	sessionID, _ := cmd.Flags().GetString(cFlag.SessionID)
	if sessionID == "" {
		input := coreSession.ReadInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = session.IDUnknown
	}

	path := nudge.PauseMarkerPath(sessionID)
	counter.Write(path, 0)
	writePause.Confirmed(cmd, sessionID)
	return nil
}
