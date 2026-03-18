//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/session"
	ctxcontext "github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// StateDir returns the project-scoped runtime state directory
// (.context/state/). Ensures the directory exists on each call — MkdirAll
// is a no-op when the directory is already present.
//
// Returns:
//   - string: Absolute path to the state directory
func StateDir() string {
	d := filepath.Join(rc.ContextDir(), dir.State)
	_ = os.MkdirAll(d, 0o750)
	return d
}

// ReadCounter reads an integer counter from a file. Returns 0 if the file
// does not exist or cannot be parsed.
//
// Parameters:
//   - path: Absolute path to the counter file
//
// Returns:
//   - int: Counter value, or 0 on error
func ReadCounter(path string) int {
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return 0
	}
	n, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
	if parseErr != nil {
		return 0
	}
	return n
}

// WriteCounter writes an integer counter to a file.
//
// Parameters:
//   - path: Absolute path to the counter file
//   - n: Counter value to write
func WriteCounter(path string, n int) {
	_ = os.WriteFile(path, []byte(strconv.Itoa(n)), 0o600)
}

// LogMessage appends a timestamped log line to the given file.
// Rotates the log when it exceeds config.HookLogMaxBytes, keeping one
// previous generation (.1 suffix) — same pattern as eventlog.
//
// Parameters:
//   - logFile: Absolute path to the log file
//   - sessionID: Session identifier (truncated to 8 chars)
//   - msg: Log message to append
func LogMessage(logFile, sessionID, msg string) {
	dir := filepath.Dir(logFile)
	_ = os.MkdirAll(dir, 0o750)

	RotateLog(logFile)

	short := sessionID
	if len(short) > 8 {
		short = short[:8]
	}

	line := fmt.Sprintf("[%s] [session:%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05"), short, msg)

	f, openErr := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec // logFile is constructed internally
	if openErr != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.WriteString(line)
}

// RotateLog checks the log file size and rotates if it exceeds
// config.HookLogMaxBytes. The previous generation is replaced.
//
// Parameters:
//   - logFile: Absolute path to the log file
func RotateLog(logFile string) {
	info, statErr := os.Stat(logFile)
	if statErr != nil {
		return
	}
	if info.Size() < int64(event.HookLogMaxBytes) {
		return
	}
	prev := logFile + ".1"
	_ = os.Remove(prev)
	_ = os.Rename(logFile, prev)
}

// IsDailyThrottled checks if a marker file was touched today (used to
// limit certain checks to once per day).
//
// Parameters:
//   - markerPath: Absolute path to the throttle marker file
//
// Returns:
//   - bool: True if the marker was touched today
func IsDailyThrottled(markerPath string) bool {
	info, statErr := os.Stat(markerPath)
	if statErr != nil {
		return false
	}
	y1, m1, d1 := info.ModTime().Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// TouchFile creates or updates the modification time of a file.
//
// Parameters:
//   - path: Absolute path to the file
func TouchFile(path string) {
	_ = os.WriteFile(path, nil, 0o600)
}

// Initialized reports whether the context directory has been properly set up
// via "ctx init". Hooks should no-op when this returns false to avoid
// creating partial state (e.g. logs/) before initialization.
//
// Returns:
//   - bool: True if context directory is initialized
func Initialized() bool {
	return ctxcontext.Initialized(rc.ContextDir())
}

// PauseMarkerPath returns the path to the session pause marker file.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - string: Absolute path to the pause marker file
func PauseMarkerPath(sessionID string) string {
	return filepath.Join(StateDir(), "ctx-paused-"+sessionID)
}

// Paused checks if the session is paused. If paused, increments the
// turn counter and returns the current count. Returns 0 if not paused.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - int: Turn count if paused, 0 if not paused
func Paused(sessionID string) int {
	path := PauseMarkerPath(sessionID)
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return 0
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	count++
	WriteCounter(path, count)
	return count
}

// PausedMessage returns the appropriate pause indicator for the given
// turn count, or empty string if not paused (turns == 0).
//
// Parameters:
//   - turns: Number of paused turns
//
// Returns:
//   - string: Pause message, or empty string
func PausedMessage(turns int) string {
	if turns == 0 {
		return ""
	}
	if turns <= 5 {
		return "ctx:paused"
	}
	return fmt.Sprintf("ctx:paused (%d turns) — resume with /ctx-resume", turns)
}

// Pause creates the session pause marker. Exported for use by the
// top-level ctx pause command.
//
// Parameters:
//   - sessionID: Session identifier
func Pause(sessionID string) {
	WriteCounter(PauseMarkerPath(sessionID), 0)
}

// Resume removes the session pause marker. Exported for use by the
// top-level ctx resume command. No-op if not paused.
//
// Parameters:
//   - sessionID: Session identifier
func Resume(sessionID string) {
	_ = os.Remove(PauseMarkerPath(sessionID))
}

// SessionStats holds the fields written to the per-session stats JSONL file.
type SessionStats struct {
	Timestamp  string `json:"ts"`
	Prompt     int    `json:"prompt"`
	Tokens     int    `json:"tokens"`
	Pct        int    `json:"pct"`
	WindowSize int    `json:"window"`
	Model      string `json:"model,omitempty"`
	Event      string `json:"event"`
}

// WriteSessionStats appends a JSONL line to .context/state/stats-{sessionID}.jsonl.
// The file is designed for `tail -f` monitoring of token usage across prompts.
// Best-effort: errors are silently ignored.
//
// Parameters:
//   - sessionID: Session identifier
//   - stats: Stats entry to write
func WriteSessionStats(sessionID string, stats SessionStats) {
	path := filepath.Join(StateDir(), "stats-"+sessionID+".jsonl")
	data, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		return
	}
	data = append(data, '\n')

	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec // state dir path
	if openErr != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(data)
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
		return session.IDUnknown
	}
	return input.SessionID
}
