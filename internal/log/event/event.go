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

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Append writes a single event to the log file.
//
// # Log-First Principle
//
// The event log is the authoritative record of "what this hook did".
// Any hook path that emits an observable side effect (webhook, stdout
// marker, state mutation) must call Append FIRST and gate the side
// effect on the log landing. If the log write fails, the side effect
// must not fire: claiming success for an event we never recorded is
// the kind of silent drift this function used to produce before it
// returned an error. See docs/security/reporting.md →
// "Log-First Audit Trail" for the rationale and call-site pattern.
//
// Noop (nil) when event logging is disabled in .ctxrc. Creates the
// state directory if it does not exist. Rotates the log when it
// exceeds EventLogMaxBytes.
//
// Parameters:
//   - event: Event type (e.g., "relay", "nudge")
//   - message: Human-readable description
//   - sessionID: Claude session ID (may be empty)
//   - detail: Optional template reference (may be nil)
//
// Returns:
//   - error: non-nil on path resolution, state-dir creation, rotation,
//     marshal, or append failure. Callers are expected to propagate
//     this error and skip any downstream webhook / state / stdout
//     side effects that would pretend the event happened. The Getwd
//     failure path is the one intentional exception: it falls back
//     to [project.FallbackName] and only warns to stderr, because
//     the event itself is still recorded, just with a less specific
//     project field. A missing CWD is never a reason to drop an
//     event entry.
func Append(
	event, message, sessionID string,
	detail *entity.TemplateRef,
) error {
	if !rc.EventLog() {
		return nil
	}

	logPath, pathErr := logFilePath()
	if pathErr != nil {
		return pathErr
	}

	// Ensure state directory exists.
	stateDir := filepath.Dir(logPath)
	if mkErr := io.SafeMkdirAll(stateDir, fs.PermExec); mkErr != nil {
		return mkErr
	}

	// Check rotation before appending.
	if rotateErr := rotate(logPath); rotateErr != nil {
		return rotateErr
	}

	projectName := project.FallbackName
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		projectName = filepath.Base(cwd)
	} else {
		// Documented fallback: record the event with a generic
		// project name rather than dropping the entry entirely.
		logWarn.Warn(warn.Getwd, cwdErr)
	}

	payload := entity.NewNotifyPayload(
		event, message, sessionID, projectName, detail,
	)

	line, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return marshalErr
	}
	newline := token.NewlineLF[0]
	line = append(line, newline)

	return io.AppendBytes(logPath, line, fs.PermFile)
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
//   - []entity.NotifyPayload: Matching events (newest last)
//   - error: Non-nil only if the log file exists but cannot be opened
func Query(opts entity.EventQueryOpts) ([]entity.NotifyPayload, error) {
	var allEvents []entity.NotifyPayload

	// Read the rotated file first (older events) if requested.
	if opts.IncludeRotated {
		prev, prevErr := prevLogFilePath()
		if prevErr != nil {
			return []entity.NotifyPayload{}, nil
		}
		events, readErr := readLogFile(prev)
		if readErr != nil {
			return nil, readErr
		}
		allEvents = append(allEvents, events...)
	}

	// Read current log file.
	current, currentErr := logFilePath()
	if currentErr != nil {
		return []entity.NotifyPayload{}, nil
	}
	events, readErr := readLogFile(current)
	if readErr != nil {
		return nil, readErr
	}
	allEvents = append(allEvents, events...)

	// Apply filters.
	var filtered []entity.NotifyPayload
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
		filtered = []entity.NotifyPayload{}
	}

	return filtered, nil
}
