//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"sort"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// SortByReadOrder sorts context files according to [config.ReadOrder].
//
// Files not in the read-order list are assigned a fallback priority
// (len(ReadOrder)) and will appear at the end. The original slice is not modified; a new sorted
// slice is returned.
//
// Parameters:
//   - files: Context files to sort
//
// Returns:
//   - []entity.FileInfo: New slice with files sorted by read priority
func SortByReadOrder(files []entity.FileInfo) []entity.FileInfo {
	// Create a map for a quick priority lookup.
	// Unknown files get fallback priority (after all known files).
	fallback := len(ctx.ReadOrder)
	priority := make(map[string]int)
	for i, name := range ctx.ReadOrder {
		priority[name] = i
	}

	// Copy and sort
	sorted := make([]entity.FileInfo, len(files))
	copy(sorted, files)

	sort.Slice(sorted, func(i, j int) bool {
		pi, ok := priority[sorted[i].Name]
		if !ok {
			pi = fallback
		}
		pj, ok := priority[sorted[j].Name]
		if !ok {
			pj = fallback
		}
		return pi < pj
	})

	return sorted
}
