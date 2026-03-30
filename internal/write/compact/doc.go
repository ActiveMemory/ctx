//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compact provides terminal output for the context compaction
// workflow (ctx compact).
//
// Functions cover the full lifecycle of a compact operation:
// [ReportHeading] opens the report, [InfoMovingTask] and
// [InfoSkippingTask] narrate per-task decisions, [InfoArchivedTasks]
// confirms what was written, [SectionsRemoved] reports empty section
// cleanup, and [ReportSummary]/[ReportClean] close the report.
//
// Example:
//
//	write.ReportHeading(cmd)
//	for _, t := range tasks {
//	    write.InfoMovingTask(cmd, t)
//	}
//	write.ReportSummary(cmd, changes)
package compact
