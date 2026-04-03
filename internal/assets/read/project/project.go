//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// Readme reads a project directory README template by directory name.
//
// Templates are stored as project/<dir>-README.md in the embedded filesystem.
//
// Parameters:
//   - dir: Directory name (e.g., "specs", "ideas")
//
// Returns:
//   - []byte: README.md content for the directory
//   - error: Non-nil if the file is not found or read fails
func Readme(dir string) ([]byte, error) {
	return assets.FS.ReadFile(
		path.Join(asset.DirProject, path.Base(dir)+asset.SuffixReadme),
	)
}
