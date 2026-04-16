//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync provides output formatting for the
// steering sync command.
//
// Steering files are tool-specific instruction files
// that guide AI agents. The sync command copies these
// files from the context directory to each tool's
// native format. This package formats the results.
//
// # Report Rendering
//
// [PrintReport] takes a steering.SyncReport and
// renders it to the command output stream. The report
// contains three categories:
//
//   - Written: files that were successfully synced.
//     Each is printed as a confirmation line.
//   - Skipped: files that already existed or were
//     unchanged. Each is printed as a skip notice.
//   - Errors: files that failed to sync. Each error
//     message is printed as a warning.
//
// After listing individual files, a summary line shows
// the count of written, skipped, and errored files.
//
// The function delegates all output to the
// write/steering package, keeping formatting logic
// separate from the sync business logic.
package sync
