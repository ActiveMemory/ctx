//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// resolvedJournalDir returns the path to the journal directory within the
// configured context directory. Uses rc.ContextDir() so it respects .ctxrc
// and CLI overrides.
func resolvedJournalDir() string {
	return filepath.Join(rc.ContextDir(), config.DirJournal)
}

// stateDir returns the project-scoped runtime state directory
// (.context/state/). Ensures the directory exists on each call — MkdirAll
// is a no-op when the directory is already present.
func stateDir() string {
	dir := filepath.Join(rc.ContextDir(), config.DirState)
	_ = os.MkdirAll(dir, 0o750)
	return dir
}

// readCounter reads an integer counter from a file. Returns 0 if the file
// does not exist or cannot be parsed.
func readCounter(path string) int {
	data, err := os.ReadFile(path) //nolint:gosec // state dir path
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return n
}

// writeCounter writes an integer counter to a file.
func writeCounter(path string, n int) {
	_ = os.WriteFile(path, []byte(strconv.Itoa(n)), 0o600)
}

// logMessage appends a timestamped log line to the given file.
// Rotates the log when it exceeds config.LogMaxBytes, keeping one
// previous generation (.1 suffix) — same pattern as eventlog.
func logMessage(logFile, sessionID, msg string) {
	dir := filepath.Dir(logFile)
	_ = os.MkdirAll(dir, 0o750)

	rotateLog(logFile)

	short := sessionID
	if len(short) > 8 {
		short = short[:8]
	}

	line := fmt.Sprintf("[%s] [session:%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05"), short, msg)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec // logFile is constructed internally
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.WriteString(line)
}

// rotateLog checks the log file size and rotates if it exceeds
// config.LogMaxBytes. The previous generation is replaced.
func rotateLog(logFile string) {
	info, statErr := os.Stat(logFile)
	if statErr != nil {
		return
	}
	if info.Size() < int64(config.LogMaxBytes) {
		return
	}
	prev := logFile + ".1"
	_ = os.Remove(prev)
	_ = os.Rename(logFile, prev)
}

// isDailyThrottled checks if a marker file was touched today (used to
// limit certain checks to once per day).
func isDailyThrottled(markerPath string) bool {
	info, err := os.Stat(markerPath)
	if err != nil {
		return false
	}
	y1, m1, d1 := info.ModTime().Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// touchFile creates or updates the modification time of a file.
func touchFile(path string) {
	_ = os.WriteFile(path, nil, 0o600)
}

// isInitialized reports whether the context directory has been properly set up
// via "ctx init". Hooks should no-op when this returns false to avoid
// creating partial state (e.g. logs/) before initialization.
func isInitialized() bool {
	return config.Initialized(rc.ContextDir())
}

// pauseMarkerPath returns the path to the session pause marker file.
func pauseMarkerPath(sessionID string) string {
	return filepath.Join(stateDir(), "ctx-paused-"+sessionID)
}

// paused checks if the session is paused. If paused, increments the
// turn counter and returns the current count. Returns 0 if not paused.
func paused(sessionID string) int {
	path := pauseMarkerPath(sessionID)
	data, readErr := os.ReadFile(path) //nolint:gosec // state dir path
	if readErr != nil {
		return 0
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	count++
	writeCounter(path, count)
	return count
}

// pausedMessage returns the appropriate pause indicator for the given
// turn count, or empty string if not paused (turns == 0).
func pausedMessage(turns int) string {
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
func Pause(sessionID string) {
	writeCounter(pauseMarkerPath(sessionID), 0)
}

// Resume removes the session pause marker. Exported for use by the
// top-level ctx resume command. No-op if not paused.
func Resume(sessionID string) {
	_ = os.Remove(pauseMarkerPath(sessionID))
}

// sessionStats holds the fields written to the per-session stats JSONL file.
type sessionStats struct {
	Timestamp  string `json:"ts"`
	Prompt     int    `json:"prompt"`
	Tokens     int    `json:"tokens"`
	Pct        int    `json:"pct"`
	WindowSize int    `json:"window"`
	Model      string `json:"model,omitempty"`
	Event      string `json:"event"`
}

// writeSessionStats appends a JSONL line to .context/state/stats-{sessionID}.jsonl.
// The file is designed for `tail -f` monitoring of token usage across prompts.
// Best-effort: errors are silently ignored.
func writeSessionStats(sessionID string, stats sessionStats) {
	path := filepath.Join(stateDir(), "stats-"+sessionID+".jsonl")
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
func ReadSessionID(stdin *os.File) string {
	input := readInput(stdin)
	if input.SessionID == "" {
		return sessionUnknown
	}
	return input.SessionID
}

// contextDirLine returns a one-line context directory identifier.
// Returns empty string if directory cannot be resolved (callers omit footer).
func contextDirLine() string {
	dir := rc.ContextDir()
	if dir == "" {
		return ""
	}
	return "Context: " + dir
}
