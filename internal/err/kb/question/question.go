//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package question

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadFile wraps a file-read failure on the
// outstanding-questions artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbQuestionReadFile), cause)
}

// MkdirDir wraps a directory-create failure on the
// outstanding-questions artifact's parent directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbQuestionMkdirDir), cause)
}

// OpenFile wraps an open-for-append failure on the
// outstanding-questions artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func OpenFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbQuestionOpenFile), cause)
}

// WriteRow wraps a row-write failure to the
// outstanding-questions artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteRow(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbQuestionAppendRow), cause)
}

// ParseQNumber wraps a strconv.Atoi failure while parsing the
// numeric portion of a Q-### token.
//
// Parameters:
//   - digits: the raw digit string that failed to parse.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ParseQNumber(digits string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbQuestionParseQNumber),
		digits, cause,
	)
}
