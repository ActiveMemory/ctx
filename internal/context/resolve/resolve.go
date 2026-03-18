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

// ResolvedJournalDir returns the path to the journal directory within the
// configured context directory.
func ResolvedJournalDir() string {
	return filepath.Join(rc.ContextDir(), dir.Journal)
}

// DirLine returns a one-line context directory identifier.
// Returns an empty string if the directory cannot be resolved.
func DirLine() string {
	d := rc.ContextDir()
	if d == "" {
		return ""
	}
	return "Context: " + d
}

// AppendDir appends a bracketed context directory footer to msg
// if a context directory is available. Returns msg unchanged otherwise.
func AppendDir(msg string) string {
	if line := DirLine(); line != "" {
		return msg + " [" + line + "]"
	}
	return msg
}
