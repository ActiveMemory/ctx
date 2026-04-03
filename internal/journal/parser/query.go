//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/entity"
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
//   - []*entity.Session: Deduplicated, filtered sessions sorted by start time
//   - error: Currently always nil (errors are silently ignored)
func findSessionsWithFilter(
	filter func(*entity.Session) bool, additionalDirs ...string,
) ([]*entity.Session, error) {
	var allSessions []*entity.Session
	scannedDirs := make(map[string]bool)

	// scanOnce scans a directory only if it hasn't been scanned yet.
	scanOnce := func(dirPath string) {
		resolved, symlinkErr := filepath.EvalSymlinks(dirPath)
		if symlinkErr != nil {
			resolved = filepath.Clean(dirPath)
		}
		if scannedDirs[resolved] {
			return
		}
		if info, statErr := os.Stat(resolved); statErr == nil && info.IsDir() {
			scannedDirs[resolved] = true
			sessions, _ := ScanDirectory(resolved)
			allSessions = append(allSessions, sessions...)
		}
	}

	// Check Claude Code default location
	home, homeErr := os.UserHomeDir()
	if homeErr == nil {
		scanOnce(filepath.Join(home, dir.Claude, dir.Projects))
	}

	// Check Copilot Chat session directories (Code + Code Insiders)
	for _, sessionDir := range CopilotSessionDirs() {
		scanOnce(sessionDir)
	}

	// Check Copilot CLI session directories (~/.copilot/ or $COPILOT_HOME)
	for _, sessionDir := range CopilotCLISessionDirs() {
		scanOnce(sessionDir)
	}

	// Check .context/sessions/ in the current working directory
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		scanOnce(filepath.Join(cwd, dir.Context, dir.Sessions))
	}

	// Check additional directories
	for _, sessionDir := range additionalDirs {
		scanOnce(sessionDir)
	}

	// Apply filter if provided
	var filtered []*entity.Session
	for _, s := range allSessions {
		if filter == nil || filter(s) {
			filtered = append(filtered, s)
		}
	}

	// Deduplicate by session ID
	seen := make(map[string]bool)
	var unique []*entity.Session
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
