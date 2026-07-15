//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcemap

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadFile wraps a stat / read failure on the source-map
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbSourcemapReadFile), cause)
}

// MkdirDir wraps a directory-create failure on the source-map
// artifact's parent directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbSourcemapMkdirDir), cause)
}

// OpenFile wraps an open-for-append failure on the source-map
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func OpenFile(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbSourcemapOpenFile), cause)
}

// WriteRow wraps a row-write failure to the source-map
// artifact.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteRow(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbSourcemapAppendRow), cause)
}
