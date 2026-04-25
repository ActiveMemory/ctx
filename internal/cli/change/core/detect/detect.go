//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package detect

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// FromMarkers finds the second most recent ctx-loaded-* marker file.
// The most recent is the current session's marker, so the reference
// point for change detection is the one before it.
//
// Returns:
//   - time.Time: Marker file modification time on success.
//   - error: [errCtx.ErrDirNotDeclared] when no context dir is
//     declared; the underlying error from [os.ReadDir] when the state
//     directory cannot be read; [os.ErrNotExist] when fewer than two
//     marker files exist (no previous session to compare against).
//     Callers treat any non-nil error as "try the next source".
func FromMarkers() (time.Time, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return time.Time{}, err
	}
	stateDir := filepath.Join(ctxDir, dir.State)
	entries, readDirErr := os.ReadDir(stateDir)
	if readDirErr != nil {
		return time.Time{}, readDirErr
	}

	type markerInfo struct {
		modTime time.Time
	}

	var markers []markerInfo
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), load_gate.PrefixCtxLoaded) {
			continue
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		markers = append(markers, markerInfo{modTime: info.ModTime()})
	}

	if len(markers) < 2 {
		// No previous-session marker on disk yet.
		return time.Time{}, os.ErrNotExist
	}

	// Sort by modtime descending.
	sort.Slice(markers, func(i, j int) bool {
		return markers[i].modTime.After(markers[j].modTime)
	})

	// Second most recent = previous session.
	return markers[1].modTime, nil
}

// FromEvents scans events.jsonl in reverse for the last
// context-load-gate event.
//
// Returns:
//   - time.Time: Event timestamp on success.
//   - error: [errCtx.ErrDirNotDeclared] when no context dir is
//     declared; the underlying error from the event log reader when
//     the file cannot be read; [os.ErrNotExist] when no matching
//     load-gate event is present or its timestamp cannot be parsed.
//     Callers treat any non-nil error as "try the next source".
func FromEvents() (time.Time, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return time.Time{}, err
	}
	eventsPath := filepath.Join(ctxDir, dir.State, event.FileLog)
	data, readErr := io.SafeReadUserFile(eventsPath)
	if readErr != nil {
		return time.Time{}, readErr
	}

	lines := strings.Split(strings.TrimSpace(string(data)), token.NewlineLF)
	// Scan in reverse for the last context-load-gate event.
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if !strings.Contains(line, load_gate.EventContextLoadGate) {
			continue
		}
		if t, ok := ExtractTimestamp(line); ok {
			return t, nil
		}
	}

	// No matching load-gate event in the log.
	return time.Time{}, os.ErrNotExist
}
