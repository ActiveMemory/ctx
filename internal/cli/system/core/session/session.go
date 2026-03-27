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
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	cfgStats "github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
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
// Both cases are harmless - hooks degrade gracefully with zero input.
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
	case <-time.After(hook.StdinReadTimeout * time.Second):
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

// LatestSessionPct returns the most recent context window usage percentage
// from the session stats JSONL. Returns 0 if no stats are available.
// This allows other hooks (e.g., check-persistence) to gate their nudges
// based on actual context window usage without re-reading the JSONL.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - int: Latest context window usage percentage (0-100), or 0 if unknown
func LatestSessionPct(sessionID string) int {
	path := filepath.Join(
		state.StateDir(),
		cfgStats.FilePrefix+sessionID+file.ExtJSONL,
	)
	data, readErr := internalIo.SafeReadUserFile(path)
	if readErr != nil {
		return 0
	}

	// Scan from the end for the last non-empty line.
	lines := splitLines(data)
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if len(line) == 0 {
			continue
		}
		var s entity.Stats
		if jsonErr := json.Unmarshal(line, &s); jsonErr != nil {
			continue
		}
		return s.Pct
	}
	return 0
}

// splitLines splits data on newline bytes, returning non-empty slices.
func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == token.NewlineLF[0] {
			if i > start {
				lines = append(lines, data[start:i])
			}
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

// WriteSessionStats appends a JSONL line to .context/state/stats-{sessionID}.jsonl.
// The file is designed for `tail -f` monitoring of token usage across prompts.
// Best-effort: errors are silently ignored.
//
// Parameters:
//   - sessionID: Session identifier
//   - stats: Stats entry to write
func WriteSessionStats(sessionID string, stats entity.Stats) {
	path := filepath.Join(
		state.StateDir(),
		cfgStats.FilePrefix+sessionID+file.ExtJSONL,
	)
	data, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		return
	}
	data = append(data, token.NewlineLF[0])

	f, openErr := internalIo.SafeAppendFile(path, fs.PermSecret)
	if openErr != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(data)
}
