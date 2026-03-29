//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// Preamble reads hook input, resolves the session ID, and checks the
// pause state. Most hooks share this exact preamble sequence.
//
// Parameters:
//   - stdin: standard input for hook JSON
//
// Returns:
//   - input: parsed hook input
//   - sessionID: resolved session identifier
//     (falls back to config.IDSessionUnknown)
//   - paused: true if the session is currently paused
func Preamble(stdin *os.File) (
	input entity.HookInput, sessionID string, paused bool,
) {
	input = coreSession.ReadInput(stdin)
	sessionID = input.SessionID
	if sessionID == "" {
		sessionID = cfgSession.IDUnknown
	}
	paused = nudge.Paused(sessionID) > 0
	return
}
