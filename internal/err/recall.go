//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// EventLogRead wraps a failure to read the event log.
//
// Parameters:
//   - cause: the underlying error from the query operation.
//
// Returns:
//   - error: "reading event log: <cause>"
func EventLogRead(cause error) error {
	return fmt.Errorf("reading event log: %w", cause)
}

// StatsGlob wraps a failure to glob stats files.
//
// Parameters:
//   - cause: the underlying error from the glob operation.
//
// Returns:
//   - error: "globbing stats files: <cause>"
func StatsGlob(cause error) error {
	return fmt.Errorf("globbing stats files: %w", cause)
}

// ReindexFileNotFound returns an error when the file to reindex does not exist.
//
// Parameters:
//   - fileName: Display name (e.g., "DECISIONS.md")
//
// Returns:
//   - error: "<fileName> not found. Run 'ctx init' first"
func ReindexFileNotFound(fileName string) error {
	return fmt.Errorf("%s not found. Run 'ctx init' first", fileName)
}

// ReindexFileRead wraps a read failure during reindexing.
//
// Parameters:
//   - filePath: Path that could not be read
//   - cause: The underlying read error
//
// Returns:
//   - error: "failed to read <filePath>: <cause>"
func ReindexFileRead(filePath string, cause error) error {
	return fmt.Errorf("failed to read %s: %w", filePath, cause)
}

// ReindexFileWrite wraps a write failure during reindexing.
//
// Parameters:
//   - filePath: Path that could not be written
//   - cause: The underlying write error
//
// Returns:
//   - error: "failed to write <filePath>: <cause>"
func ReindexFileWrite(filePath string, cause error) error {
	return fmt.Errorf("failed to write %s: %w", filePath, cause)
}

// OpenLogFile wraps a failure to open a log file.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to open log file: <cause>"
func OpenLogFile(cause error) error {
	return fmt.Errorf("failed to open log file: %w", cause)
}
