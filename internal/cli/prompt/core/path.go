//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// PromptsDir returns the path to the prompts directory.
//
// Returns:
//   - string: Absolute path to .context/prompts/
func PromptsDir() string {
	return filepath.Join(rc.ContextDir(), dir.Prompts)
}
