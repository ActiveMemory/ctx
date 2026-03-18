//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ReadFile wraps a session file read failure.
//
// Parameters:
//   - cause: the underlying error from reading the file.
//
// Returns:
//   - error: "read file: <cause>"
func ReadFile(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserReadFile), cause,
	)
}

// OpenFile wraps a session file open failure.
//
// Parameters:
//   - cause: the underlying error from opening the file.
//
// Returns:
//   - error: "open file: <cause>"
func OpenFile(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserOpenFile), cause,
	)
}

// NoMatch returns an error when no parser can handle a file.
//
// Parameters:
//   - path: the file path that no parser matched.
//
// Returns:
//   - error: "no parser found for file: <path>"
func NoMatch(path string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserNoMatch), path,
	)
}

// WalkDir wraps a directory walk failure during session scanning.
//
// Parameters:
//   - cause: the underlying error from filepath.Walk.
//
// Returns:
//   - error: "walk directory: <cause>"
func WalkDir(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserWalkDir), cause,
	)
}

// FileError wraps a per-file parse failure with the file path.
//
// Parameters:
//   - path: the file path that failed to parse.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "<path>: <cause>"
func FileError(path string, cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserFileError), path, cause,
	)
}

// ScanFile wraps a session file scan failure.
//
// Parameters:
//   - cause: the underlying error from scanning the file.
//
// Returns:
//   - error: "scan file: <cause>"
func ScanFile(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserScanFile), cause,
	)
}

// Unmarshal wraps a JSON unmarshal failure during session parsing.
//
// Parameters:
//   - cause: the underlying error from JSON unmarshaling.
//
// Returns:
//   - error: "unmarshal: <cause>"
func Unmarshal(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserUnmarshal), cause,
	)
}

// ParseFile wraps a failure to parse a file.
//
// Parameters:
//   - path: file path that could not be parsed
//   - cause: the underlying parse error
//
// Returns:
//   - error: "failed to parse <path>: <cause>"
func ParseFile(path string, cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidationParseFile), path, cause,
	)
}
