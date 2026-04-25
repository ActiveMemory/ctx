//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/write/session"
)

// Run executes the resume command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: Session ID from flag (empty to read from stdin)
//
// Returns:
//   - error: Always nil
func Run(cmd *cobra.Command, sessionID string) error {
	if sessionID == "" {
		sessionID = coreSession.ReadID(os.Stdin)
	}
	if resumeErr := nudge.Resume(sessionID); resumeErr != nil {
		return resumeErr
	}
	session.Resumed(cmd, sessionID)
	return nil
}
