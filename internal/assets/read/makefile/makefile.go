//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package makefile

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// Ctx reads the ctx-owned Makefile include template.
//
// Returns:
//   - []byte: Makefile.ctx content
//   - error: Non-nil if the file is not found or read fails
func Ctx() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathMakefileCtx)
}
