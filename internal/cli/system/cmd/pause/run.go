//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"os"

	"github.com/spf13/cobra"

	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
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
	sessionID, ok := coreCheck.PausePreamble(cmd, stdin)
	if !ok {
		return nil
	}
	if pauseErr := nudge.Pause(sessionID); pauseErr != nil {
		logWarn.Warn(warn.StateDirProbe, pauseErr)
		return nil
	}
	writePause.Confirmed(cmd, sessionID)
	return nil
}
