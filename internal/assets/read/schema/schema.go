//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// Schema reads the embedded JSON Schema for .ctxrc.
//
// Returns:
//   - []byte: JSON Schema content
//   - error: Non-nil if the file is not found or read fails
func Schema() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCtxrcSchema)
}
