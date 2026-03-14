//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	time2 "github.com/ActiveMemory/ctx/internal/config/time"
)

// UUIDPattern matches a UUID (v4) anywhere in a filename.
var UUIDPattern = regexp.MustCompile(
	`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
)

// AutoPrune silently removes session-scoped state files older than the
// given number of days. Called from context-load-gate on session start.
// Returns the number of files removed. Errors are swallowed — auto-prune
// is best-effort and must never block session startup.
//
// Parameters:
//   - days: Prune files older than this many days
//
// Returns:
//   - int: Number of files pruned
func AutoPrune(days int) int {
	dir := StateDir()

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return 0
	}

	cutoff := time.Now().Add(-time.Duration(days) * time2.HoursPerDay * time.Hour)
	var pruned int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !UUIDPattern.MatchString(entry.Name()) {
			continue
		}

		info, statErr := entry.Info()
		if statErr != nil {
			continue
		}

		if info.ModTime().After(cutoff) {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if rmErr := os.Remove(path); rmErr == nil {
			pruned++
		}
	}

	return pruned
}

// FormatAge formats a time.Time as a human-readable age string.
//
// Parameters:
//   - t: Time to format
//
// Returns:
//   - string: Age string (e.g. "5m", "3h", "2d")
func FormatAge(t time.Time) string {
	d := time.Since(t)
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
