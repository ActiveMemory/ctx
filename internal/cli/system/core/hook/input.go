//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
)

// FormatContext builds a JSON HookResponse with additionalContext for the
// given hook event. This is the standard way for non-blocking hooks to inject
// directives that the agent will actually process (plain text gets ignored).
//
// Parameters:
//   - event: Hook event name
//   - context: Additional context string
//
// Returns:
//   - string: JSON-encoded hook response
func FormatContext(event, context string) string {
	resp := HookResponse{
		HookSpecificOutput: &HookSpecificOutput{
			HookEventName:     event,
			AdditionalContext: context,
		},
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

// Preamble reads hook input, resolves the session ID, and checks the
// pause state. Most hooks share this exact preamble sequence.
//
// Parameters:
//   - stdin: standard input for hook JSON
//
// Returns:
//   - input: parsed hook input
//   - sessionID: resolved session identifier (falls back to config.IDSessionUnknown)
//   - paused: true if the session is currently paused
func Preamble(stdin *os.File) (
	input coreSession.HookInput, sessionID string, paused bool,
) {
	input = coreSession.ReadInput(stdin)
	sessionID = input.SessionID
	if sessionID == "" {
		sessionID = cfgSession.IDUnknown
	}
	paused = core.Paused(sessionID) > 0
	return
}
