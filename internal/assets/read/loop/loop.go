//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// RalphTemplate reads a Ralph-mode template file by name.
//
// Ralph mode templates are designed for autonomous loop operation,
// with instructions for one-task-per-iteration, completion signals,
// and no clarifying questions.
//
// Parameters:
//   - name: Template filename (e.g., "PROMPT.md")
//
// Returns:
//   - []byte: Template content from ralph/
//   - error: Non-nil if the file is not found or read fails
func RalphTemplate(name string) ([]byte, error) {
	return assets.FS.ReadFile(path.Join(asset.DirRalph, name))
}
