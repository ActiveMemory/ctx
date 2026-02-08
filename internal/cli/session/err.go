//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import "fmt"

// errIndexNotPositive returns an error when a parsed index is less than 1.
//
// Returns:
//   - error: "index must be positive"
func errIndexNotPositive() error {
	return fmt.Errorf("index must be positive")
}

// errSectionNotFound returns an error when a header section is missing.
//
// Parameters:
//   - header: Section header that was not found
//
// Returns:
//   - error: Formatted error with the missing header name
func errSectionNotFound(header string) error {
	return fmt.Errorf("section not found: %s", header)
}

// errReadingFile wraps a file read error with context.
//
// Parameters:
//   - err: Underlying read error
//
// Returns:
//   - error: Wrapped error with format "error reading file: <cause>"
func errReadingFile(err error) error {
	return fmt.Errorf("error reading file: %w", err)
}

// errNoSessionsDir returns an error when the sessions directory does not exist.
//
// Returns:
//   - error: Formatted error suggesting "ctx session save"
func errNoSessionsDir() error {
	return fmt.Errorf(
		"no sessions directory found. Run 'ctx session save' first",
	)
}

// errReadSession wraps a failure to read a session file.
//
// Parameters:
//   - err: Underlying read error
//
// Returns:
//   - error: Wrapped error with format "failed to read session file: <cause>"
func errReadSession(err error) error {
	return fmt.Errorf("failed to read session file: %w", err)
}

// errFileNotFound returns an error when an input file does not exist.
//
// Parameters:
//   - path: Path that was not found
//
// Returns:
//   - error: Formatted error with the missing path
func errFileNotFound(path string) error {
	return fmt.Errorf("file not found: %s", path)
}

// errExtractInsights wraps a failure to extract insights from a transcript.
//
// Parameters:
//   - err: Underlying extraction error
//
// Returns:
//   - error: Wrapped error with format "failed to extract insights: <cause>"
func errExtractInsights(err error) error {
	return fmt.Errorf("failed to extract insights: %w", err)
}

// errParseTranscript wraps a failure to parse a JSONL transcript.
//
// Parameters:
//   - err: Underlying parse error
//
// Returns:
//   - error: Wrapped error with format "failed to parse transcript: <cause>"
func errParseTranscript(err error) error {
	return fmt.Errorf("failed to parse transcript: %w", err)
}

// errWriteOutput wraps a failure to write output to a file.
//
// Parameters:
//   - err: Underlying write error
//
// Returns:
//   - error: Wrapped error with format "failed to write output: <cause>"
func errWriteOutput(err error) error {
	return fmt.Errorf("failed to write output: %w", err)
}

// errCreateSessionsDir wraps a failure to create the sessions directory.
//
// Parameters:
//   - err: Underlying mkdir error
//
// Returns:
//   - error: Wrapped error with format "failed to create sessions directory: <cause>"
func errCreateSessionsDir(err error) error {
	return fmt.Errorf("failed to create sessions directory: %w", err)
}

// errBuildContent wraps a failure to build session content.
//
// Parameters:
//   - err: Underlying build error
//
// Returns:
//   - error: Wrapped error with format "failed to build session content: <cause>"
func errBuildContent(err error) error {
	return fmt.Errorf("failed to build session content: %w", err)
}

// errWriteSession wraps a failure to write a session file.
//
// Parameters:
//   - err: Underlying write error
//
// Returns:
//   - error: Wrapped error with format "failed to write session file: <cause>"
func errWriteSession(err error) error {
	return fmt.Errorf("failed to write session file: %w", err)
}

// errNoSessions returns an error when no session files exist.
//
// Returns:
//   - error: "no sessions found"
func errNoSessions() error {
	return fmt.Errorf("no sessions found")
}

// errIndexOutOfRange returns an error when a session index exceeds the count.
//
// Parameters:
//   - idx: Requested index
//   - count: Total number of sessions
//
// Returns:
//   - error: Formatted error with the valid range
func errIndexOutOfRange(idx, count int) error {
	return fmt.Errorf("index %d out of range (1-%d)", idx, count)
}

// errNoSessionMatch returns an error when no session matches a query.
//
// Parameters:
//   - query: Search query that found no matches
//
// Returns:
//   - error: Formatted error with the query
func errNoSessionMatch(query string) error {
	return fmt.Errorf("no session found matching %q", query)
}

// errMultipleMatches returns an error when a query matches multiple sessions.
//
// Parameters:
//   - query: Search query
//   - matches: List of matching filenames
//
// Returns:
//   - error: Formatted error listing the matches
func errMultipleMatches(query string, matches []string) error {
	return fmt.Errorf("multiple sessions match %q: %v", query, matches)
}

// errReadSessionsDir wraps a failure to read the sessions directory.
//
// Parameters:
//   - err: Underlying read error
//
// Returns:
//   - error: Wrapped error with format "failed to read sessions directory: <cause>"
func errReadSessionsDir(err error) error {
	return fmt.Errorf("failed to read sessions directory: %w", err)
}
