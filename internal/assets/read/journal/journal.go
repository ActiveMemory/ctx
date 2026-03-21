//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// ExtraCSS reads the embedded extra.css for journal site generation.
//
// Returns:
//   - []byte: CSS content
//   - error: Non-nil if the file is not found or read fails
func ExtraCSS() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathExtraCSS)
}
