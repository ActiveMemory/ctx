//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"os"

	"github.com/spf13/cobra"

	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	writeSession "github.com/ActiveMemory/ctx/internal/write/session"
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
	sessionID, ok := coreCheck.PausePreamble(cmd, stdin)
	if !ok {
		return nil
	}
	if resumeErr := nudge.Resume(sessionID); resumeErr != nil {
		ctxLog.Warn(warn.StateDirProbe, resumeErr)
		return nil
	}
	writeSession.Resumed(cmd, sessionID)
	return nil
}
