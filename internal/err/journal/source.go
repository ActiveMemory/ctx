//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// EventLogRead wraps a failure to read the event log.
//
// Parameters:
//   - cause: the underlying error from the query operation.
//
// Returns:
//   - error: "reading event log: <cause>"
func EventLogRead(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrJournalSourceEventLogRead), cause,
	)
}

// StatsGlob wraps a failure to glob stats files.
//
// Parameters:
//   - cause: the underlying error from the glob operation.
//
// Returns:
//   - error: "globbing stats files: <cause>"
func StatsGlob(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrJournalSourceStatsGlob), cause,
	)
}

// OpenLogFile wraps a failure to open a log file.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to open log file: <cause>"
func OpenLogFile(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrJournalSourceOpenLogFile), cause,
	)
}
