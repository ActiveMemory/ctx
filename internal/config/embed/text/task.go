//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for task archive output.
const (
	DescKeyTaskArchiveDryRunBlock    = "task-archive.dry-run-block"
	DescKeyTaskArchiveNoCompleted    = "task-archive.no-completed"
	DescKeyTaskArchivePendingRemain  = "task-archive.pending-remain"
	DescKeyTaskArchiveSkipIncomplete = "task-archive.skip-incomplete"
	DescKeyTaskArchiveSkipping       = "task-archive.skipping"
	DescKeyTaskArchiveSuccess        = "task-archive.success"
	DescKeyTaskArchiveSuccessWithAge = "task-archive.success-with-age"
)

// DescKeys for task snapshot output.
const (
	DescKeyTaskSnapshotHeaderFormat  = "task-snapshot.header-format"
	DescKeyTaskSnapshotCreatedFormat = "task-snapshot.created-format"
	DescKeyTaskSnapshotSaved         = "task-snapshot.saved"
)

// DescKeys for task completion check nudge.
const (
	DescKeyCheckTaskCompletionFallback     = "check-task-completion.fallback"
	DescKeyCheckTaskCompletionNudgeMessage = "check-task-completion.nudge-message"
)

// DescKeys for task management write output.
const (
	DescKeyWriteCompletedTask = "write.completed-task"
	DescKeyWriteMovingTask    = "write.moving-task"
)
