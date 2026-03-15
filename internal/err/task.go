//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"
)

// TaskFileNotFound returns an error when TASKS.md does not exist.
//
// Returns:
//   - error: "TASKS.md not found"
func TaskFileNotFound() error {
	return fmt.Errorf("TASKS.md not found")
}

// TaskFileRead wraps a failure to read TASKS.md.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "failed to read TASKS.md: <cause>"
func TaskFileRead(cause error) error {
	return fmt.Errorf("failed to read TASKS.md: %w", cause)
}

// TaskFileWrite wraps a failure to write TASKS.md.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "failed to write TASKS.md: <cause>"
func TaskFileWrite(cause error) error {
	return fmt.Errorf("failed to write TASKS.md: %w", cause)
}

// TaskMultipleMatches returns an error when a query matches more than one task.
//
// Parameters:
//   - query: the search string that matched multiple tasks
//
// Returns:
//   - error: "multiple tasks match <query>; be more specific or use task number"
func TaskMultipleMatches(query string) error {
	return fmt.Errorf(
		"multiple tasks match %q; be more specific or use task number",
		query,
	)
}

// TaskNotFound returns an error when no task matches the query.
//
// Parameters:
//   - query: the search string that matched nothing
//
// Returns:
//   - error: "no task matching <query> found"
func TaskNotFound(query string) error {
	return fmt.Errorf("no task matching %q found", query)
}

// NoCompletedTasks returns an error when there are no completed tasks to archive.
//
// Returns:
//   - error: "no completed tasks to archive"
func NoCompletedTasks() error {
	return fmt.Errorf("no completed tasks to archive")
}

// NoTaskSpecified returns an error when no task query was provided.
//
// Returns:
//   - error: "no task specified"
func NoTaskSpecified() error {
	return errors.New("no task specified")
}

// NoTaskMatch returns an error when no task matches the search query.
//
// Parameters:
//   - query: the search string that matched nothing.
//
// Returns:
//   - error: "no task matching \"<query>\" found"
func NoTaskMatch(query string) error {
	return fmt.Errorf("no task matching %q found", query)
}

// SnapshotWrite wraps a failure to write a task snapshot file.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to write snapshot: <cause>"
func SnapshotWrite(cause error) error {
	return fmt.Errorf("failed to write snapshot: %w", cause)
}
