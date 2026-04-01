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
)

// EssentialFilesPresent reports whether contextDir contains at least one of
// the essential context files (TASKS.md, CONSTITUTION.md, DECISIONS.md). A
// directory with only logs/ or other non-essential content is considered
// uninitialized.
//
// Parameters:
//   - contextDir: Absolute path to the context directory to inspect
//
// Returns:
//   - bool: True if at least one essential file exists
func EssentialFilesPresent(contextDir string) bool {
	for _, f := range ctx.FilesRequired {
		if _, statErr := os.Stat(filepath.Join(contextDir, f)); statErr == nil {
			return true
		}
	}
	return false
}
