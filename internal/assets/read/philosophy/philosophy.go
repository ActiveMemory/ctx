//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package philosophy

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/file"
)

// WhyDoc reads a "why" document by name from the embedded filesystem.
//
// Parameters:
//   - name: Document name (e.g., "manifesto", "about", "design-invariants")
//
// Returns:
//   - []byte: Document content from why/
//   - error: Non-nil if the file is not found or read fails
func WhyDoc(name string) ([]byte, error) {
	return assets.FS.ReadFile(path.Join(asset.DirWhy, name+file.ExtMarkdown))
}

// WhyDocList returns available "why" document names (without extension).
//
// Returns:
//   - []string: List of document names in why/
//   - error: Non-nil if directory read fails
func WhyDocList() ([]string, error) {
	entries, readErr := assets.FS.ReadDir(asset.DirWhy)
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if len(name) > 3 && name[len(name)-3:] == file.ExtMarkdown {
				names = append(names, name[:len(name)-3])
			}
		}
	}
	return names, nil
}
