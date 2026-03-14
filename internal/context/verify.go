//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Exists checks if a context directory exists.
//
// If dir is empty, it uses the configured context directory.
//
// Parameters:
//   - dir: Directory path to check, or empty string for default
//
// Returns:
//   - bool: True if the directory exists and is a directory
//
// Initialized reports whether the context directory contains all required files.
//
// Parameters:
//   - contextDir: Directory path to check
//
// Returns:
//   - bool: True if all required context files exist
func Initialized(contextDir string) bool {
	for _, f := range ctx.FilesRequired {
		if _, err := os.Stat(filepath.Join(contextDir, f)); err != nil {
			return false
		}
	}
	return true
}

func Exists(dir string) bool {
	if dir == "" {
		dir = rc.ContextDir()
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

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
