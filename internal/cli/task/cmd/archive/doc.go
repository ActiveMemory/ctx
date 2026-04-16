//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package archive implements the "ctx task archive"
// cobra subcommand.
//
// This command moves completed tasks from TASKS.md to
// a timestamped archive file in .context/archive/.
// Pending tasks remain in TASKS.md so active work is
// never lost.
//
// # Usage
//
//	ctx task archive [--dry-run]
//
// # Flags
//
//	--dry-run   Preview which tasks would be archived
//	            and what the remaining TASKS.md would
//	            look like, without modifying any files.
//
// # Behavior
//
// The command delegates to task/core/archive in two
// phases:
//
//  1. Plan: scans TASKS.md for completed tasks (those
//     marked with [x]). Tasks whose children include
//     incomplete items are skipped to avoid orphaning
//     pending work. The plan returns the archivable
//     set, skipped names, and remaining pending count.
//
//  2. Execute: writes the archivable tasks to a new
//     timestamped file in .context/archive/ and
//     rewrites TASKS.md with only the pending tasks.
//
// When --dry-run is set only the plan phase runs and
// the preview is printed.
//
// # Output
//
// Without --dry-run: prints the count of archived
// tasks, the archive file path, and the remaining
// pending count. With --dry-run: prints a preview of
// the archive content separated from the remaining
// tasks.
//
// When no completed tasks exist, prints a message
// indicating nothing to archive. When tasks are
// skipped due to incomplete children, lists the
// skipped task names.
//
// # Delegation
//
// Planning and execution are in task/core/archive.
// Output formatting is in write/archive.
package archive
