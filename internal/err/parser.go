//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ParserReadFile wraps a session file read failure.
//
// Parameters:
//   - cause: the underlying error from reading the file.
//
// Returns:
//   - error: "read file: <cause>"
func ParserReadFile(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserReadFile), cause,
	)
}

// ParserOpenFile wraps a session file open failure.
//
// Parameters:
//   - cause: the underlying error from opening the file.
//
// Returns:
//   - error: "open file: <cause>"
func ParserOpenFile(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserOpenFile), cause,
	)
}

// ParserNoMatch returns an error when no parser can handle a file.
//
// Parameters:
//   - path: the file path that no parser matched.
//
// Returns:
//   - error: "no parser found for file: <path>"
func ParserNoMatch(path string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserNoMatch), path,
	)
}

// ParserWalkDir wraps a directory walk failure during session scanning.
//
// Parameters:
//   - cause: the underlying error from filepath.Walk.
//
// Returns:
//   - error: "walk directory: <cause>"
func ParserWalkDir(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserWalkDir), cause,
	)
}

// ParserFileError wraps a per-file parse failure with the file path.
//
// Parameters:
//   - path: the file path that failed to parse.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "<path>: <cause>"
func ParserFileError(path string, cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserFileError), path, cause,
	)
}

// ParserScanFile wraps a session file scan failure.
//
// Parameters:
//   - cause: the underlying error from scanning the file.
//
// Returns:
//   - error: "scan file: <cause>"
func ParserScanFile(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserScanFile), cause,
	)
}

// ParserUnmarshal wraps a JSON unmarshal failure during session parsing.
//
// Parameters:
//   - cause: the underlying error from JSON unmarshaling.
//
// Returns:
//   - error: "unmarshal: <cause>"
func ParserUnmarshal(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrParserUnmarshal), cause,
	)
}
