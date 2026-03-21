//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"path/filepath"

	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// ReadOrder returns context file paths in the recommended reading order.
//
// Files are ordered according to [config.ReadOrder] and filtered to
// exclude empty files. Paths are returned as full paths relative to the
// context directory.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: File paths in reading order (e.g., ".context/CONSTITUTION.md")
func ReadOrder(ctx *entity.Context) []string {
	var order []string
	for _, name := range ctxCfg.ReadOrder {
		if f := ctx.File(name); f != nil && !f.IsEmpty {
			order = append(order, filepath.Join(ctx.Dir, f.Name))
		}
	}
	return order
}
