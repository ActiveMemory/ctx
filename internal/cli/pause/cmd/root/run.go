//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Run executes the pause command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: Session ID from flag (empty to read from stdin)
//
// Returns:
//   - error: Always nil
func Run(cmd *cobra.Command, sessionID string) error {
	if sessionID == "" {
		sessionID = core.ReadSessionID(os.Stdin)
	}
	core.Pause(sessionID)
	cmd.Println(fmt.Sprintf("Context hooks paused for session %s", sessionID))
	return nil
}
