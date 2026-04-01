//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Append writes a single event to the log file.
//
// Noop when event logging is disabled in .ctxrc. Creates the state
// directory if it does not exist. Rotates the log when it exceeds
// EventLogMaxBytes. All errors are silently ignored: event logging
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
	if mkErr := os.MkdirAll(stateDir, fs.PermExec); mkErr != nil {
		return
	}

	// Check rotation before appending.
	rotate(logPath)

	projectName := project.FallbackName
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		projectName = filepath.Base(cwd)
	} else {
		logWarn.Warn(warn.Getwd, cwdErr)
	}

	payload := notify.Payload{
		Event:     event,
		Message:   message,
		Detail:    detail,
		SessionID: sessionID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   projectName,
	}

	line, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return
	}
	newline := token.NewlineLF[0]
	line = append(line, newline)

	io.AppendBytes(logPath, line, fs.PermFile)
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
func Query(opts entity.EventQueryOpts) ([]notify.Payload, error) {
	var allEvents []notify.Payload

	// Read the rotated file first (older events) if requested.
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
