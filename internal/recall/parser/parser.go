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
	"strings"
)

// registeredParsers holds all available session parsers.
// Add new parsers here when supporting additional tools.
var registeredParsers = []SessionParser{
	NewClaudeCodeParser(),
	NewMarkdownSessionParser(),
}

// ParseFile parses a session file using the appropriate parser.
//
// It auto-detects the file format by trying each registered parser.
//
// Parameters:
//   - path: Path to the session file to parse
//
// Returns:
//   - []*Session: All sessions found in the file
//   - error: Non-nil if no parser can handle the file or parsing fails
func ParseFile(path string) ([]*Session, error) {
	for _, parser := range registeredParsers {
		if parser.Matches(path) {
			return parser.ParseFile(path)
		}
	}
	return nil, fmt.Errorf("no parser found for file: %s", path)
}

// ScanDirectory recursively scans a directory for session files.
//
// It finds all parseable files, parses them, and aggregates sessions.
// Sessions are sorted by start time (newest first). Parse errors for
// individual files are silently ignored; use ScanDirectoryWithErrors
// if you need to report them.
//
// Parameters:
//   - dir: Root directory to scan recursively
//
// Returns:
//   - []*Session: All sessions found, sorted by start time (newest first)
//   - error: Non-nil if directory traversal fails
func ScanDirectory(dir string) ([]*Session, error) {
	sessions, _, err := ScanDirectoryWithErrors(dir)
	return sessions, err
}

// ScanDirectoryWithErrors is like ScanDirectory but also returns parse errors.
//
// Use this when you want to report files that failed to parse while still
// returning successfully parsed sessions.
//
// Parameters:
//   - dir: Root directory to scan recursively
//
// Returns:
//   - []*Session: Successfully parsed sessions, sorted by start time
//   - []error: Errors from files that failed to parse
//   - error: Non-nil if directory traversal fails
func ScanDirectoryWithErrors(dir string) ([]*Session, []error, error) {
	var allSessions []*Session
	var parseErrors []error

	err := filepath.Walk(dir, func(
		path string, info os.FileInfo, err error,
	) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip subagents directories - they contain sidechain sessions
			// that share the parent sessionId and would cause duplicates
			if info.Name() == "subagents" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip files in paths containing /subagents/ (defensive check)
		if strings.Contains(path, string(filepath.Separator)+"subagents"+string(filepath.Separator)) {
			return nil
		}

		// Try to parse with any registered parser
		for _, parser := range registeredParsers {
			if parser.Matches(path) {
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
// Parameters:
//   - additionalDirs: Optional additional directories to scan
//
// Returns:
//   - []*Session: Deduplicated sessions sorted by start time (newest first)
//   - error: Non-nil if scanning fails (partial results may still be returned)
func FindSessions(additionalDirs ...string) ([]*Session, error) {
	return findSessionsWithFilter(nil, additionalDirs...)
}

// FindSessionsForCWD searches for sessions matching the given
// working directory.
//
// Matching is done in order of preference:
//  1. Git remote URL match - if both directories are git repos with
//     the same remote
//  2. Path relative to home - e.g., "WORKSPACE/ctx" matches across users
//  3. Exact CWD match - fallback for non-git, non-home paths
//
// Parameters:
//   - cwd: Working directory to filter by
//   - additionalDirs: Optional additional directories to scan
//
// Returns:
//   - []*Session: Filtered sessions sorted by start time (newest first)
//   - error: Non-nil if scanning fails
func FindSessionsForCWD(
	cwd string, additionalDirs ...string,
) ([]*Session, error) {
	// Get current project's git remote (if available)
	currentRemote := gitRemote(cwd)

	// Get path relative to home directory
	currentRelPath := getPathRelativeToHome(cwd)

	return findSessionsWithFilter(func(s *Session) bool {
		// 1. Try git remote match (most robust)
		if currentRemote != "" {
			sessionRemote := gitRemote(s.CWD)
			if sessionRemote != "" && sessionRemote == currentRemote {
				return true
			}
		}

		// 2. Try the path relative to the home match
		if currentRelPath != "" {
			sessionRelPath := getPathRelativeToHome(s.CWD)
			if sessionRelPath != "" && sessionRelPath == currentRelPath {
				return true
			}
		}

		// 3. Fallback to an exact match
		return s.CWD == cwd
	}, additionalDirs...)
}

// Parser returns a parser for the specified tool.
//
// Parameters:
//   - tool: Tool identifier (e.g., "claude-code")
//
// Returns:
//   - SessionParser: The parser for the tool, or nil if not found
func Parser(tool string) SessionParser {
	for _, parser := range registeredParsers {
		if parser.Tool() == tool {
			return parser
		}
	}
	return nil
}

// RegisteredTools returns the list of supported tools.
//
// Returns:
//   - []string: Tool identifiers for all registered parsers
func RegisteredTools() []string {
	tools := make([]string, len(registeredParsers))
	for i, parser := range registeredParsers {
		tools[i] = parser.Tool()
	}
	return tools
}
