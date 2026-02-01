//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"os"

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
func Exists(dir string) bool {
	if dir == "" {
		dir = rc.GetContextDir()
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}
