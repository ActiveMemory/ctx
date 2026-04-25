//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"errors"
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
		if _, statErr := os.Stat(filepath.Join(contextDir, f)); statErr != nil {
			return false
		}
	}
	return true
}

// Exists checks whether a context directory exists.
//
// If dir is empty, it uses the configured context directory. A missing
// context declaration is a "false, nil" result (no configured dir, so
// nothing to check); a resolver failure for any other reason and a stat
// failure that is not "does not exist" are both propagated so callers
// can distinguish "the directory is not there" from "we could not find
// out."
//
// Parameters:
//   - dir: path to check, or empty string for default
//
// Returns:
//   - bool: true if the path exists and is a directory
//   - error: non-nil on resolver failure (other than not-declared) or
//     stat failure (other than not-exist)
func Exists(dir string) (bool, error) {
	if dir == "" {
		declared, err := rc.ContextDir()
		if err != nil {
			return false, err
		}
		dir = declared
	}
	info, statErr := os.Stat(dir)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return false, nil
		}
		return false, statErr
	}
	return info.IsDir(), nil
}
