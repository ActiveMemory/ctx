//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sort orders context files by priority and
// recency for status display.
//
// # Priority Sorting
//
// [FilesByPriority] sorts a file slice in place using
// the priority order defined in rc.FilePriority.
// CONSTITUTION.md sorts first, followed by TASKS.md,
// CONVENTIONS.md, and other context files. This ensures
// the status display presents files in the recommended
// reading order for agents and users.
//
// # Recency Sorting
//
// [RecentFiles] returns the n most recently modified
// files. It copies the input slice to avoid mutation,
// sorts by modification time descending, and truncates
// to n entries. The status command uses this to show a
// recent activity section highlighting files that
// changed most recently.
//
// Both functions operate on []entity.FileInfo slices
// from the loaded context.
package sort
