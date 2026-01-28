//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// registeredParsers holds all available session parsers.
// Add new parsers here when supporting additional tools.
var registeredParsers = []SessionParser{
	NewClaudeCodeParser(),
}

// ParseFile parses a session file using the appropriate parser.
//
// It auto-detects the file format by trying each registered parser.
// Returns an error if no parser can handle the file.
func ParseFile(path string) ([]*Session, error) {
	for _, parser := range registeredParsers {
		if parser.CanParse(path) {
			return parser.ParseFile(path)
		}
	}
	return nil, fmt.Errorf("no parser found for file: %s", path)
}

// ScanDirectory recursively scans a directory for session files.
//
// It finds all parseable files, parses them, and aggregates sessions.
// Sessions are sorted by start time (newest first).
func ScanDirectory(dir string) ([]*Session, error) {
	var allSessions []*Session
	var parseErrors []error

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Try to parse with any registered parser
		for _, parser := range registeredParsers {
			if parser.CanParse(path) {
				sessions, err := parser.ParseFile(path)
				if err != nil {
					parseErrors = append(parseErrors, fmt.Errorf("%s: %w", path, err))
					break
				}
				allSessions = append(allSessions, sessions...)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}

	// Sort by start time (newest first)
	sort.Slice(allSessions, func(i, j int) bool {
		return allSessions[i].StartTime.After(allSessions[j].StartTime)
	})

	return allSessions, nil
}

// ScanDirectoryWithErrors is like ScanDirectory but also returns parse errors.
//
// Use this when you want to report files that failed to parse while still
// returning successfully parsed sessions.
func ScanDirectoryWithErrors(dir string) ([]*Session, []error, error) {
	var allSessions []*Session
	var parseErrors []error

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Try to parse with any registered parser
		for _, parser := range registeredParsers {
			if parser.CanParse(path) {
				sessions, err := parser.ParseFile(path)
				if err != nil {
					parseErrors = append(parseErrors, fmt.Errorf("%s: %w", path, err))
					break
				}
				allSessions = append(allSessions, sessions...)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("walk directory: %w", err)
	}

	// Sort by start time (newest first)
	sort.Slice(allSessions, func(i, j int) bool {
		return allSessions[i].StartTime.After(allSessions[j].StartTime)
	})

	return allSessions, parseErrors, nil
}

// FindSessions searches for session files in common locations.
//
// It checks:
//  1. ~/.claude/projects/ (Claude Code default)
//  2. The specified directory (if provided)
//
// Returns all found sessions sorted by start time.
func FindSessions(additionalDirs ...string) ([]*Session, error) {
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

	// Check additional directories
	for _, dir := range additionalDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			sessions, _ := ScanDirectory(dir)
			allSessions = append(allSessions, sessions...)
		}
	}

	// Deduplicate by session ID
	seen := make(map[string]bool)
	var unique []*Session
	for _, s := range allSessions {
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

// GetParser returns a parser for the specified tool.
//
// Returns nil if no parser is registered for the tool.
func GetParser(tool string) SessionParser {
	for _, parser := range registeredParsers {
		if parser.Tool() == tool {
			return parser
		}
	}
	return nil
}

// RegisteredTools returns the list of supported tools.
func RegisteredTools() []string {
	tools := make([]string, len(registeredParsers))
	for i, parser := range registeredParsers {
		tools[i] = parser.Tool()
	}
	return tools
}
