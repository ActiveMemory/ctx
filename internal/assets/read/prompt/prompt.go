//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// TemplateList returns available prompt template file names.
//
// Returns:
//   - []string: List of template filenames in prompt-templates/
//   - error: Non-nil if directory read fails
func TemplateList() ([]string, error) {
	entries, err := assets.FS.ReadDir(asset.DirPromptTemplates)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// Template reads a prompt template by name.
//
// Parameters:
//   - name: Template filename (e.g., "code-review.md")
//
// Returns:
//   - []byte: Template content from prompt-templates/
//   - error: Non-nil if the file is not found or read fails
func Template(name string) ([]byte, error) {
	return assets.FS.ReadFile(path.Join(asset.DirPromptTemplates, name))
}
