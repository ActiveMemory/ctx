//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/ActiveMemory/ctx/internal/config"
)

// findSessionsWithFilter scans common locations and additional directories
// for session files, applying an optional filter.
//
// It checks ~/.claude/projects/ (Claude Code default) and any additional
// directories provided. Results are deduplicated by session ID and sorted
// by start time (newest first).
//
// Parameters:
//   - filter: Optional function to filter sessions (nil includes all)
//   - additionalDirs: Optional additional directories to scan
//
// Returns:
//   - []*Session: Deduplicated, filtered sessions sorted by start time
//   - error: Currently always nil (errors are silently ignored)
func findSessionsWithFilter(
	filter func(*Session) bool, additionalDirs ...string,
) ([]*Session, error) {
	var allSessions []*Session

	// Check Claude Code default location
	home, err := os.UserHomeDir()
	if err == nil {
		claudeDir := filepath.Join(home, ".claude", "projects")
		if info, err := os.Stat(claudeDir); err == nil && info.IsDir() {
			sessions, _ := ScanDirectory(claudeDir)
			allSessions = append(allSessions, sessions...)
		}
	}

	// Check .context/sessions/ in the current working directory
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		sessionsDir := filepath.Join(cwd, config.DirContext, config.DirSessions)
		if info, statErr := os.Stat(sessionsDir); statErr == nil && info.IsDir() {
			sessions, _ := ScanDirectory(sessionsDir)
			allSessions = append(allSessions, sessions...)
		}
	}

	// Check additional directories
	for _, dir := range additionalDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			sessions, _ := ScanDirectory(dir)
			allSessions = append(allSessions, sessions...)
		}
	}

	// Apply filter if provided
	var filtered []*Session
	for _, s := range allSessions {
		if filter == nil || filter(s) {
			filtered = append(filtered, s)
		}
	}

	// Deduplicate by session ID
	seen := make(map[string]bool)
	var unique []*Session
	for _, s := range filtered {
		if !seen[s.ID] {
			seen[s.ID] = true
			unique = append(unique, s)
		}
	}

	// Sort by start time (newest first)
	sort.Slice(unique, func(i, j int) bool {
		return unique[i].StartTime.After(unique[j].StartTime)
	})

	return unique, nil
}
