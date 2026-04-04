//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for task archive output.
const (
	// DescKeyTaskArchiveDryRunBlock is the text key for task archive dry run
	// block messages.
	DescKeyTaskArchiveDryRunBlock = "task-archive.dry-run-block"
	// DescKeyTaskArchiveNoCompleted is the text key for task archive no completed
	// messages.
	DescKeyTaskArchiveNoCompleted = "task-archive.no-completed"
	// DescKeyTaskArchivePendingRemain is the text key for task archive pending
	// remain messages.
	DescKeyTaskArchivePendingRemain = "task-archive.pending-remain"
	// DescKeyTaskArchiveSkipIncomplete is the text key for task archive skip
	// incomplete messages.
	DescKeyTaskArchiveSkipIncomplete = "task-archive.skip-incomplete"
	// DescKeyTaskArchiveSkipping is the text key for task archive skipping
	// messages.
	DescKeyTaskArchiveSkipping = "task-archive.skipping"
	// DescKeyTaskArchiveSuccess is the text key for task archive success messages.
	DescKeyTaskArchiveSuccess = "task-archive.success"
	// DescKeyTaskArchiveSuccessWithAge is the text key for task archive success
	// with age messages.
	DescKeyTaskArchiveSuccessWithAge = "task-archive.success-with-age"
)

// DescKeys for task snapshot output.
const (
	// DescKeyTaskSnapshotHeaderFormat is the text key for task snapshot header
	// format messages.
	DescKeyTaskSnapshotHeaderFormat = "task-snapshot.header-format"
	// DescKeyTaskSnapshotCreatedFormat is the text key for task snapshot created
	// format messages.
	DescKeyTaskSnapshotCreatedFormat = "task-snapshot.created-format"
	// DescKeyTaskSnapshotSaved is the text key for task snapshot saved messages.
	DescKeyTaskSnapshotSaved = "task-snapshot.saved"
)

// DescKeys for task completion check nudge.
const (
	// DescKeyCheckTaskCompletionFallback is the text key for check task
	// completion fallback messages.
	DescKeyCheckTaskCompletionFallback = "check-task-completion.fallback"
	// DescKeyCheckTaskCompletionNudgeMessage is the text key for check task
	// completion nudge message messages.
	DescKeyCheckTaskCompletionNudgeMessage = "check-task-completion.nudge-message"
)

// DescKeys for task management write output.
const (
	// DescKeyWriteCompletedTask is the text key for write completed task messages.
	DescKeyWriteCompletedTask = "write.completed-task"
	// DescKeyWriteMovingTask is the text key for write moving task messages.
	DescKeyWriteMovingTask = "write.moving-task"
)
