//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// FormatContext builds a JSON Response with additionalContext for the
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
	resp := response{
		Output: &responseOutput{
			HookEventName:     event,
			AdditionalContext: context,
		},
	}
	data, _ := json.Marshal(resp)
	return string(data)
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
//   - entity.HookInput: Parsed input or zero value
func ReadInput(r io.Reader) entity.HookInput {
	if f, ok := r.(*os.File); ok {
		if fi, readErr := f.Stat(); readErr == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return entity.HookInput{}
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

	var input entity.HookInput
	select {
	case res := <-ch:
		if res.err == nil {
			_ = json.Unmarshal(res.data, &input)
		}
	case <-time.After(2 * time.Second):
	}
	return input
}

// ReadSessionID reads the session ID from stdin JSON, returning the
// fallback "unknown" if stdin is empty or unparseable.
//
// Parameters:
//   - stdin: File to read input from
//
// Returns:
//   - string: Session ID or config.IDSessionUnknown
func ReadSessionID(stdin *os.File) string {
	input := ReadInput(stdin)
	if input.SessionID == "" {
		return cfgSession.IDUnknown
	}
	return input.SessionID
}

// WriteSessionStats appends a JSONL line to .context/state/stats-{sessionID}.jsonl.
// The file is designed for `tail -f` monitoring of token usage across prompts.
// Best-effort: errors are silently ignored.
//
// Parameters:
//   - sessionID: Session identifier
//   - stats: Stats entry to write
func WriteSessionStats(sessionID string, stats entity.Stats) {
	path := filepath.Join(state.StateDir(), "stats-"+sessionID+".jsonl")
	data, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		return
	}
	data = append(data, '\n')

	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.PermSecret) //nolint:gosec // state dir path
	if openErr != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(data)
}
