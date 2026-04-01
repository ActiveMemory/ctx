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
// The most recent is the current session's marker.
//
// Returns:
//   - time.Time: Marker file modification time
//   - bool: True if a valid marker was found
func FromMarkers() (time.Time, bool) {
	stateDir := filepath.Join(rc.ContextDir(), dir.State)
	entries, readDirErr := os.ReadDir(stateDir)
	if readDirErr != nil {
		return time.Time{}, false
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
		return time.Time{}, false
	}

	// Sort by modtime descending.
	sort.Slice(markers, func(i, j int) bool {
		return markers[i].modTime.After(markers[j].modTime)
	})

	// Second most recent = previous session.
	return markers[1].modTime, true
}

// FromEvents scans events.jsonl in reverse for the last
// context-load-gate event.
//
// Returns:
//   - time.Time: Event timestamp
//   - bool: True if a valid event was found
func FromEvents() (time.Time, bool) {
	eventsPath := filepath.Join(rc.ContextDir(), dir.State, event.FileLog)
	data, readErr := io.SafeReadUserFile(eventsPath)
	if readErr != nil {
		return time.Time{}, false
	}

	lines := strings.Split(strings.TrimSpace(string(data)), token.NewlineLF)
	// Scan in reverse for the last context-load-gate event.
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if !strings.Contains(line, load_gate.EventContextLoadGate) {
			continue
		}
		if t, ok := ExtractTimestamp(line); ok {
			return t, true
		}
	}

	return time.Time{}, false
}
