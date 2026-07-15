//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contradiction

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadFile wraps a file-read failure on the contradictions
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbContradictionReadFile), cause)
}

// MkdirDir wraps a directory-create failure on the
// contradictions artifact's parent directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbContradictionMkdirDir), cause)
}

// OpenFile wraps an open-for-append failure on the
// contradictions artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func OpenFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbContradictionOpenFile), cause)
}

// WriteRow wraps a row-write failure to the contradictions
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteRow(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbContradictionAppendRow), cause)
}

// ParseCNumber wraps a strconv.Atoi failure while parsing the
// numeric portion of a C-### token.
//
// Parameters:
//   - digits: the raw digit string that failed to parse.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ParseCNumber(digits string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbContradictionParseCNumber),
		digits, cause,
	)
}
