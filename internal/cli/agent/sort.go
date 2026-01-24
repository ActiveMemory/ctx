//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"fmt"

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
		for _, f := range ctx.Files {
			if f.Name == name && !f.IsEmpty {
				order = append(order, fmt.Sprintf("%s/%s", ctx.Dir, f.Name))
				break
			}
		}
	}
	return order
}
