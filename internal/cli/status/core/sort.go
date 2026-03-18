//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"sort"

	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// SortFilesByPriority sorts files in-place by the recommended read order.
//
// Uses rc.FilePriority to determine ordering (CONSTITUTION first,
// then TASKS, CONVENTIONS, etc.).
//
// Parameters:
//   - files: Slice of files to sort (modified in place)
func SortFilesByPriority(files []entity.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return rc.FilePriority(
			files[i].Name,
		) < rc.FilePriority(files[j].Name)
	})
}

// RecentFilesSorted returns the n most recently modified files.
//
// Parameters:
//   - files: Source files to select from
//   - n: Maximum number of files to return
//
// Returns:
//   - []entity.FileInfo: Up to n files sorted by modification time
//     (newest first)
func RecentFilesSorted(files []entity.FileInfo, n int) []entity.FileInfo {
	sorted := make([]entity.FileInfo, len(files))
	copy(sorted, files)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ModTime.After(sorted[j].ModTime)
	})
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	return sorted
}
