//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package eventlog provides append-only JSONL event logging for hook
// diagnostics. Events are written to .context/state/events.jsonl when
// enabled via event_log: true in .ctxrc. The log format is identical
// to webhook payloads (notify.Payload) — one struct, two sinks.
package eventlog

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Append writes a single event to the log file.
//
// Noop when event logging is disabled in .ctxrc. Creates the state
// directory if it does not exist. Rotates the log when it exceeds
// EventLogMaxBytes. All errors are silently ignored — event logging
// must never break hook execution.
//
// Parameters:
//   - event: Event type (e.g., "relay", "nudge")
//   - message: Human-readable description
//   - sessionID: Claude session ID (may be empty)
//   - detail: Optional template reference (may be nil)
func Append(event, message, sessionID string, detail *notify.TemplateRef) {
	if !rc.EventLog() {
		return
	}

	logPath := logFilePath()

	// Ensure state directory exists.
	stateDir := filepath.Dir(logPath)
	if mkErr := os.MkdirAll(stateDir, config.PermExec); mkErr != nil {
		return
	}

	// Check rotation before appending.
	rotate(logPath)

	project := "unknown"
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		project = filepath.Base(cwd)
	}

	payload := notify.Payload{
		Event:     event,
		Message:   message,
		Detail:    detail,
		SessionID: sessionID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   project,
	}

	line, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return
	}
	line = append(line, '\n')

	//nolint:gosec // project-local state path
	f, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, config.PermFile)
	if openErr != nil {
		return
	}
	defer f.Close() //nolint:errcheck // best-effort logging

	_, _ = f.Write(line)
}

// QueryOpts controls event filtering and pagination.
type QueryOpts struct {
	Hook           string // filter by hook name (from detail)
	Session        string // filter by session ID
	Event          string // filter by event type
	Last           int    // return last N events (0 = all)
	IncludeRotated bool   // also read events.1.jsonl
}

// Query reads events from the log, applying filters.
//
// Returns an empty slice (not nil) when no events match or the log
// file does not exist. Malformed lines are silently skipped.
//
// Parameters:
//   - opts: Filter and limit options
//
// Returns:
//   - []notify.Payload: Matching events (newest last)
//   - error: Non-nil only if the log file exists but cannot be opened
func Query(opts QueryOpts) ([]notify.Payload, error) {
	var allEvents []notify.Payload

	// Read rotated file first (older events) if requested.
	if opts.IncludeRotated {
		prev := prevLogFilePath()
		events, readErr := readLogFile(prev)
		if readErr != nil {
			return nil, readErr
		}
		allEvents = append(allEvents, events...)
	}

	// Read current log file.
	current := logFilePath()
	events, readErr := readLogFile(current)
	if readErr != nil {
		return nil, readErr
	}
	allEvents = append(allEvents, events...)

	// Apply filters.
	var filtered []notify.Payload
	for _, e := range allEvents {
		if !matchesFilter(e, opts) {
			continue
		}
		filtered = append(filtered, e)
	}

	// Apply --last limit.
	if opts.Last > 0 && len(filtered) > opts.Last {
		filtered = filtered[len(filtered)-opts.Last:]
	}

	if filtered == nil {
		filtered = []notify.Payload{}
	}

	return filtered, nil
}

// readLogFile reads and parses all events from a JSONL file.
// Returns empty slice if the file does not exist.
func readLogFile(path string) ([]notify.Payload, error) {
	f, openErr := os.Open(path) //nolint:gosec // project-local state path
	if openErr != nil {
		if os.IsNotExist(openErr) {
			return nil, nil
		}
		return nil, openErr
	}
	defer f.Close() //nolint:errcheck // read-only

	var events []notify.Payload
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var p notify.Payload
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &p); unmarshalErr != nil {
			continue // skip malformed lines
		}
		events = append(events, p)
	}

	return events, nil
}

// matchesFilter reports whether an event matches the query filters.
func matchesFilter(e notify.Payload, opts QueryOpts) bool {
	if opts.Event != "" && e.Event != opts.Event {
		return false
	}
	if opts.Session != "" && e.SessionID != opts.Session {
		return false
	}
	if opts.Hook != "" {
		if e.Detail == nil || e.Detail.Hook != opts.Hook {
			return false
		}
	}
	return true
}

// rotate checks the current log file size and rotates if needed.
// Best-effort: errors are silently ignored.
func rotate(logPath string) {
	info, statErr := os.Stat(logPath)
	if statErr != nil {
		return // file doesn't exist yet, nothing to rotate
	}
	if info.Size() < int64(config.EventLogMaxBytes) {
		return
	}

	prevPath := prevLogFilePath()
	_ = os.Remove(prevPath)
	_ = os.Rename(logPath, prevPath)
}

// logFilePath returns the path to the current event log.
func logFilePath() string {
	return filepath.Join(rc.ContextDir(), config.DirState, config.FileEventLog)
}

// prevLogFilePath returns the path to the rotated event log.
func prevLogFilePath() string {
	return filepath.Join(rc.ContextDir(), config.DirState, config.FileEventLogPrev)
}
