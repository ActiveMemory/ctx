//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package timeline

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadFile wraps a stat / read failure on the timeline
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbTimelineReadFile), cause)
}

// MkdirDir wraps a directory-create failure on the timeline
// artifact's parent directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbTimelineMkdirDir), cause)
}

// OpenFile wraps an open-for-append failure on the timeline
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func OpenFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbTimelineOpenFile), cause)
}

// WriteRow wraps a row-write failure to the timeline
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteRow(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbTimelineWriteRow), cause)
}
