//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/file"
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
	for _, f := range file.FilesRequired {
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
