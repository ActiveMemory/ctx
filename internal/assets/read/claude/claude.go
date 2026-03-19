//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides access to Claude Code integration files
// embedded in the assets filesystem.
package claude

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// Md reads the CLAUDE.md template from the embedded filesystem.
//
// CLAUDE.md is deployed to the project root by a dedicated handler
// during initialization, separate from the .context/ templates.
//
// Returns:
//   - []byte: CLAUDE.md content
//   - error: Non-nil if the file is not found or read fails
func Md() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCLAUDEMd)
}
