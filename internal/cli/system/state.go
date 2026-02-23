//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
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

// secureTempDir returns a user-specific temp directory for ctx state files.
// Uses $XDG_RUNTIME_DIR when available (tmpfs, user-owned, 0700 on Linux),
// otherwise creates a user-specific subdirectory under os.TempDir().
func secureTempDir() string {
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		dir := filepath.Join(xdg, "ctx")
		_ = os.MkdirAll(dir, 0o700)
		return dir
	}
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("ctx-%d", os.Getuid()))
	_ = os.MkdirAll(dir, 0o700)
	return dir
}

// readCounter reads an integer counter from a file. Returns 0 if the file
// does not exist or cannot be parsed.
func readCounter(path string) int {
	data, err := os.ReadFile(path) //nolint:gosec // temp file path
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
func logMessage(logFile, sessionID, msg string) {
	dir := filepath.Dir(logFile)
	_ = os.MkdirAll(dir, 0o750)

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
	dir := rc.ContextDir()
	for _, f := range config.FilesRequired {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			return false
		}
	}
	return true
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
