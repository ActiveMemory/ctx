//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/session"
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
	input HookInput, sessionID string, paused bool,
) {
	input = ReadInput(stdin)
	sessionID = input.SessionID
	if sessionID == "" {
		sessionID = session.IDUnknown
	}
	paused = core.Paused(sessionID) > 0
	return
}

// ReadInput reads and parses the JSON hook input from r.
// Returns a zero-value HookInput on any error (graceful degradation).
//
// Guards against blocking forever on stdin:
//   - Terminal (character device): returns immediately
//   - Pipe/file with no EOF within 2s: times out and returns zero value
//
// Both cases are harmless — hooks degrade gracefully with zero input.
//
// Parameters:
//   - r: Reader to read hook input from
//
// Returns:
//   - HookInput: Parsed input or zero value
func ReadInput(r io.Reader) HookInput {
	if f, ok := r.(*os.File); ok {
		if fi, readErr := f.Stat(); readErr == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return HookInput{}
		}
	}

	type readResult struct {
		data []byte
		err  error
	}
	ch := make(chan readResult, 1)
	go func() {
		data, readErr := io.ReadAll(r)
		ch <- readResult{data, readErr}
	}()

	var input HookInput
	select {
	case res := <-ch:
		if res.err == nil {
			_ = json.Unmarshal(res.data, &input)
		}
	case <-time.After(2 * time.Second):
	}
	return input
}
