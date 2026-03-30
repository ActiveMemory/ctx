//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import "github.com/ActiveMemory/ctx/internal/entity"

// Result holds the outcome of an archive operation for the caller to report.
//
// Fields:
//   - Archivable: Tasks eligible for archival
//   - SkippedNames: Task names skipped (incomplete subtasks)
//   - PendingCount: Tasks still pending
//   - Content: Formatted archive file content
//   - ArchivePath: Path to the written archive file
//   - NewTasksBody: Updated TASKS.md content with archived tasks removed
type Result struct {
	Archivable   []entity.TaskBlock
	SkippedNames []string
	PendingCount int
	Content      string
	ArchivePath  string
	NewTasksBody string
}
