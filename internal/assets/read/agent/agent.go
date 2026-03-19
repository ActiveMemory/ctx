//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent provides access to agent integration files embedded
// in the assets filesystem.
package agent

import (
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
