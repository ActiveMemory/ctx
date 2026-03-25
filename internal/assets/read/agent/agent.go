//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent provides access to agent integration files embedded
// in the assets filesystem.
package agent

import (
	"io/fs"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// CopilotInstructions reads the embedded Copilot instructions template.
//
// Returns:
//   - []byte: Template content from hooks/copilot-instructions.md
//   - error: Non-nil if the file is not found or read fails
func CopilotInstructions() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCopilotInstructions)
}

// CopilotCLIHooksJSON reads the embedded Copilot CLI hooks config.
//
// Returns:
//   - []byte: JSON content from hooks/copilot-cli/ctx-hooks.json
//   - error: Non-nil if the file is not found or read fails
func CopilotCLIHooksJSON() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCopilotCLIHooksJSON)
}

// CopilotCLIScripts reads all embedded Copilot CLI hook scripts.
// Returns a map of filename to content for scripts in
// hooks/copilot-cli/scripts/.
//
// Returns:
//   - map[string][]byte: Filename -> content for each script
//   - error: Non-nil if the directory read fails
func CopilotCLIScripts() (map[string][]byte, error) {
	scripts := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirHooksCopilotCLIScrp)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".sh") && !strings.HasSuffix(name, ".ps1") {
			continue
		}
		content, readErr := assets.FS.ReadFile(asset.DirHooksCopilotCLIScrp + "/" + name)
		if readErr != nil {
			return nil, readErr
		}
		scripts[name] = content
	}
	return scripts, nil
}
