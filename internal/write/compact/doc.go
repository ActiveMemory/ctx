//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compact provides terminal output for the
// context compaction workflow (ctx compact).
//
// # Report Lifecycle
//
// Functions cover the full lifecycle of a compact
// operation. [ReportHeading] opens the report with a
// title and separator. Per-task decisions are narrated
// by [InfoMovingTask] (completed task being archived)
// and [InfoSkippingTask] (task excluded due to
// incomplete children). [InfoArchivedTasks] confirms
// how many tasks were written to the archive file
// with the age threshold used.
//
// # Cleanup Output
//
// [SectionsRemoved] reports how many empty sections
// were pruned from a context file during compaction.
// [TaskError] reports errors encountered while
// processing individual tasks.
//
// # Summary
//
// [ReportSummary] prints the final change count, and
// [ReportClean] prints the message when no changes
// were needed.
//
// # Message Categories
//
//   - Info: per-task decisions, archive results
//   - Error: task processing failures
//   - Summary: total change count or clean status
//
// # Usage
//
//	compact.ReportHeading(cmd)
//	for _, t := range tasks {
//	    compact.InfoMovingTask(cmd, t)
//	}
//	compact.ReportSummary(cmd, changes)
package compact
