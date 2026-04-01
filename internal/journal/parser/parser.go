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
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
)

// registeredParsers holds all available session parsers.
// Add new parsers here when supporting additional tools.
var registeredParsers = []Session{
	NewClaudeCode(),
	NewCopilot(),
	NewCopilotCLI(),
	NewMarkdownSession(),
}

// ParseFile parses a session file using the appropriate parser.
//
// It auto-detects the file format by trying each registered parser.
//
// Parameters:
//   - path: Path to the session file to parse
//
// Returns:
//   - []*entity.Session: All sessions found in the file
//   - error: Non-nil if no parser can handle the file or parsing fails
func ParseFile(path string) ([]*entity.Session, error) {
	for _, p := range registeredParsers {
		if p.Matches(path) {
			return p.ParseFile(path)
		}
	}
	return nil, errParser.NoMatch(path)
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
//   - []*entity.Session: All sessions found, sorted by
//     start time (newest first)
//   - error: Non-nil if directory traversal fails
func ScanDirectory(dir string) ([]*entity.Session, error) {
	sessions, _, scanErr := ScanDirectoryWithErrors(dir)
	return sessions, scanErr
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
//   - []*entity.Session: Successfully parsed sessions, sorted by start time
//   - []error: Errors from files that failed to parse
//   - error: Non-nil if directory traversal fails
func ScanDirectoryWithErrors(dir string) ([]*entity.Session, []error, error) {
	var allSessions []*entity.Session
	var parseErrors []error

	walkErr := filepath.Walk(dir, func(
		path string, info os.FileInfo, err error,
	) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip subagents directories - they contain sidechain sessions
			// that share the parent sessionId and would cause duplicates
			if info.Name() == parser.DirSubagents {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip files in paths containing /subagents/ (defensive check)
		sep := string(filepath.Separator)
		subagentPath := sep + parser.DirSubagents + sep
		if strings.Contains(path, subagentPath) {
			return nil
		}

		// Try to parse with any registered parser
		for _, p := range registeredParsers {
			if p.Matches(path) {
				sessions, err := p.ParseFile(path)
				if err != nil {
					parseErrors = append(parseErrors, errParser.FileError(path, err))
					break
				}
				allSessions = append(allSessions, sessions...)
				break
			}
		}

		return nil
	})

	if walkErr != nil {
		return nil, nil, errParser.WalkDir(walkErr)
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
//   - []*entity.Session: Deduplicated sessions sorted by
//     start time (newest first)
//   - error: Non-nil if scanning fails (partial results may still be returned)
func FindSessions(additionalDirs ...string) ([]*entity.Session, error) {
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
//   - []*entity.Session: Filtered sessions sorted by start time (newest first)
//   - error: Non-nil if scanning fails
func FindSessionsForCWD(
	cwd string, additionalDirs ...string,
) ([]*entity.Session, error) {
	// Get current project's git remote (if available)
	currentRemote := gitRemote(cwd)

	// Get path relative to home directory
	currentRelPath := getPathRelativeToHome(cwd)

	return findSessionsWithFilter(func(s *entity.Session) bool {
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
//   - Session: The parser for the tool, or nil if not found
func Parser(tool string) Session {
	for _, p := range registeredParsers {
		if p.Tool() == tool {
			return p
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
	for i, p := range registeredParsers {
		tools[i] = p.Tool()
	}
	return tools
}
