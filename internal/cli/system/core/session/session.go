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
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/parse"
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
	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		ctxLog.Warn(warn.Marshal, marshalErr)
		return ""
	}
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
		fi, readErr := f.Stat()
		if readErr == nil && fi.Mode()&os.ModeCharDevice != 0 {
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

// ReadID reads the session ID from stdin JSON, returning the
// fallback "unknown" if stdin is empty or unparseable.
//
// Parameters:
//   - stdin: File to read input from
//
// Returns:
//   - string: Session ID or config.IDSessionUnknown
func ReadID(stdin *os.File) string {
	input := ReadInput(stdin)
	if input.SessionID == "" {
		return cfgSession.IDUnknown
	}
	return input.SessionID
}

// LatestPct returns the most recent context window usage percentage
// from the session stats JSONL. Returns 0 if no stats are available.
// This allows other hooks (e.g., check-persistence) to gate their nudges
// based on actual context window usage without re-reading the JSONL.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - int: Latest context window usage percentage (0-100), or 0 if unknown
func LatestPct(sessionID string) int {
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return 0
	}
	path := filepath.Join(
		stateDir,
		cfgStats.FilePrefix+sessionID+file.ExtJSONL,
	)
	data, readErr := internalIo.SafeReadUserFile(path)
	if readErr != nil {
		return 0
	}

	// Scan from the end for the last non-empty line.
	lines := parse.ByteLines(data)
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

// WriteStats appends a JSONL line to .context/state/stats-{sessionID}.jsonl.
// The file is designed for `tail -f` monitoring of token usage across prompts.
// Errors are propagated; see Returns for the rationale.
//
// Parameters:
//   - sessionID: Session identifier
//   - stats: Stats entry to write
//
// Returns:
//   - error: non-nil when marshaling or the append fails. Stats are
//     an audit trail of per-session token usage; surfacing a write
//     failure lets callers honour the log-first principle (do not
//     claim success for a session action whose stats entry never
//     landed).
func WriteStats(sessionID string, stats entity.Stats) error {
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return dirErr
	}
	path := filepath.Join(
		stateDir,
		cfgStats.FilePrefix+sessionID+file.ExtJSONL,
	)
	data, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		return marshalErr
	}
	data = append(data, token.NewlineLF[0])

	return internalIo.AppendBytes(path, data, fs.PermSecret)
}
