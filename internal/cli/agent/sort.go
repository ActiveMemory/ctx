//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// getReadOrder returns context file paths in the recommended reading order.
//
// Files are ordered according to [config.FileReadOrder] and filtered to
// exclude empty files. Paths are returned as full paths relative to the
// context directory.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: File paths in reading order (e.g., ".context/CONSTITUTION.md")
func getReadOrder(ctx *context.Context) []string {
	var order []string
	for _, name := range config.FileReadOrder {
		if f := ctx.File(name); f != nil && !f.IsEmpty {
			order = append(order, filepath.Join(ctx.Dir, f.Name))
		}
	}
	return order
}
