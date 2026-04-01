//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// Message reads a hook message template by hook name and filename.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//   - filename: Template filename (e.g., "gate.txt")
//
// Returns:
//   - []byte: Template content from hooks/messages/<hook>/
//   - error: Non-nil if the file is not found or read fails
func Message(hook, filename string) ([]byte, error) {
	return assets.FS.ReadFile(path.Join(asset.DirHooksMessages, hook, filename))
}

// MessageRegistry reads the embedded registry.yaml that describes
// all hook message templates.
//
// Returns:
//   - []byte: Raw YAML content
//   - error: Non-nil if the file is not found or read fails
func MessageRegistry() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathMessageRegistry)
}

// MessageList returns available hook message directory names.
//
// Each hook is a directory under hooks/messages/ containing one or
// more variant .txt template files.
//
// Returns:
//   - []string: List of hook directory names
//   - error: Non-nil if directory read fails
func MessageList() ([]string, error) {
	entries, readErr := assets.FS.ReadDir(asset.DirHooksMessages)
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// VariantList returns available variant filenames for a hook.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//
// Returns:
//   - []string: List of variant filenames (e.g., "gate.txt")
//   - error: Non-nil if the hook directory is not found or read fails
func VariantList(hook string) ([]string, error) {
	entries, readErr := assets.FS.ReadDir(path.Join(asset.DirHooksMessages, hook))
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
