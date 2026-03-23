//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tidy

import "github.com/ActiveMemory/ctx/internal/entity"

// fileUpdate holds the new content for a context file that changed
// during compaction.
type fileUpdate struct {
	Path    string
	Content []byte
}

// sectionClean records how many empty sections were removed from a file.
type sectionClean struct {
	FileName string
	Removed  int
}

// CompactResult holds the outcome of a CompactContext call.
//
// Callers decide how to report results (CLI prints, MCP returns
// JSON-RPC responses) and how to write files (os.WriteFile, etc.).
type CompactResult struct {
	// TasksMoved lists the parent text of each task moved to Completed.
	TasksMoved []string
	// TasksSkipped lists parent text of completed tasks with pending children.
	TasksSkipped []string
	// TasksFileUpdate is non-nil when TASKS.md content changed.
	TasksFileUpdate *fileUpdate
	// ArchivableBlocks are blocks eligible for archival.
	ArchivableBlocks []entity.TaskBlock
	// SectionsCleaned lists files where empty sections were removed.
	SectionsCleaned []sectionClean
	// SectionFileUpdates holds new content for files with sections removed.
	SectionFileUpdates []fileUpdate
}

// TotalChanges returns the number of items compacted.
func (r *CompactResult) TotalChanges() int {
	total := len(r.TasksMoved)
	for _, sc := range r.SectionsCleaned {
		total += sc.Removed
	}
	return total
}
