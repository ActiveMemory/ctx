//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// JournalDir returns the path to the journal directory within the
// configured context directory.
//
// Returns:
//   - string: Absolute path to the journal directory
func JournalDir() string {
	return filepath.Join(rc.ContextDir(), dir.Journal)
}

// DirLine returns a one-line context directory identifier.
// Returns an empty string if the directory cannot be resolved.
//
// Returns:
//   - string: "Context: <path>" line, or empty string if unresolved
func DirLine() string {
	d := rc.ContextDir()
	if d == "" {
		return ""
	}
	return "Context: " + d
}

// AppendDir appends a bracketed context directory footer to msg
// if a context directory is available. Returns msg unchanged otherwise.
//
// Parameters:
//   - msg: Base message to append the directory footer to
//
// Returns:
//   - string: Message with appended "[Context: <path>]", or msg unchanged
func AppendDir(msg string) string {
	if line := DirLine(); line != "" {
		return msg + " [" + line + "]"
	}
	return msg
}
