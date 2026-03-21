//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"fmt"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
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
	sessionID, _ := cmd.Flags().GetString("session-id")
	if sessionID == "" {
		input := core.ReadInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = session.IDUnknown
	}

	path := core.PauseMarkerPath(sessionID)
	core.WriteCounter(path, 0)
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyPauseConfirmed), sessionID),
	)
	return nil
}
