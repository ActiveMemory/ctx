//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// FileNotFound returns an error when TASKS.md does not exist.
//
// Returns:
//   - error: "TASKS.md not found"
func FileNotFound() error {
	return errors.New(desc.TextDesc(text.DescKeyErrTaskFileNotFound))
}

// FileRead wraps a failure to read TASKS.md.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "failed to read TASKS.md: <cause>"
func FileRead(cause error) error {
	return fmt.Errorf(desc.TextDesc(text.DescKeyErrTaskFileRead), cause)
}

// FileWrite wraps a failure to write TASKS.md.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "failed to write TASKS.md: <cause>"
func FileWrite(cause error) error {
	return fmt.Errorf(desc.TextDesc(text.DescKeyErrTaskFileWrite), cause)
}

// MultipleMatches returns an error when a query matches more than one task.
//
// Parameters:
//   - query: the search string that matched multiple tasks
//
// Returns:
//   - error: "multiple tasks match <query>; be more specific or use task number"
func MultipleMatches(query string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrTaskMultipleMatches), query,
	)
}

// NotFound returns an error when no task matches the query.
//
// Parameters:
//   - query: the search string that matched nothing
//
// Returns:
//   - error: "no task matching <query> found"
func NotFound(query string) error {
	return fmt.Errorf(desc.TextDesc(text.DescKeyErrTaskNotFound), query)
}

// NoneCompleted returns an error when there are no completed tasks to archive.
//
// Returns:
//   - error: "no completed tasks to archive"
func NoneCompleted() error {
	return errors.New(desc.TextDesc(text.DescKeyErrTaskNoCompletedTasks))
}

// NoneSpecified returns an error when no task query was provided.
//
// Returns:
//   - error: "no task specified"
func NoneSpecified() error {
	return errors.New(desc.TextDesc(text.DescKeyErrTaskNoTaskSpecified))
}

// NoMatch returns an error when no task matches the search query.
//
// Parameters:
//   - query: the search string that matched nothing.
//
// Returns:
//   - error: "no task matching \"<query>\" found"
func NoMatch(query string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrTaskNoTaskMatch), query,
	)
}

// SnapshotWrite wraps a failure to write a task snapshot file.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to write snapshot: <cause>"
func SnapshotWrite(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrTaskSnapshotWrite), cause,
	)
}
