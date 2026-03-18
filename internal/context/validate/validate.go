//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Initialized reports whether the context directory contains all
// required files.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//
// Returns:
//   - bool: true if every required context file exists
func Initialized(contextDir string) bool {
	for _, f := range ctx.FilesRequired {
		if _, err := os.Stat(filepath.Join(contextDir, f)); err != nil {
			return false
		}
	}
	return true
}

// Exists checks whether a context directory exists.
//
// If dir is empty, it uses the configured context directory.
//
// Parameters:
//   - dir: path to check, or empty string for default
//
// Returns:
//   - bool: true if the path exists and is a directory
func Exists(dir string) bool {
	if dir == "" {
		dir = rc.ContextDir()
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}
