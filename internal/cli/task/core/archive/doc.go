//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package archive moves completed tasks from TASKS.md
// into timestamped archive files under .context/archive/.
// It uses a plan-then-execute pattern so the cmd/ layer
// can show a dry-run preview before writing anything.
//
// # Planning Phase
//
// [Plan] reads TASKS.md, parses it into task blocks via
// the tidy package, and classifies each block as
// archivable or skipped. A block is archivable when all
// of its subtasks are complete. The returned [Result]
// contains:
//
//   - Archivable blocks ready to move
//   - Skipped task names (have incomplete subtasks)
//   - The pending task count after archival
//   - Pre-formatted archive file content
//   - The new TASKS.md body with archived blocks removed
//
// Plan does not write any files, making it safe for
// dry-run reporting.
//
// # Execution Phase
//
// [Execute] takes a Result from Plan, writes the archive
// content to a new timestamped file via tidy.WriteArchive,
// and overwrites TASKS.md with the trimmed body. Returns
// the path to the created archive file.
//
// # Data Flow
//
//  1. cmd/ layer calls Plan to get the Result
//  2. cmd/ layer shows dry-run preview or confirms
//  3. cmd/ layer calls Execute with the Result
//  4. Archive file written, TASKS.md updated
package archive
