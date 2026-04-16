//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides business logic for the status
// command, which displays context directory health and
// file information.
//
// The core package is a namespace that groups three
// subpackages: out, preview, and sort. Together they
// transform a loaded context entity into formatted
// output for the user.
//
// # Subpackages
//
// The out subpackage renders the full status display in
// either JSON or text format. [out.PersistStatusJSON]
// encodes context metadata and file status as indented
// JSON. [out.PersistStatusText] renders a formatted
// text report with file indicators, token counts, and
// recent activity.
//
// The preview subpackage extracts short content
// previews from context files. [preview.Content]
// returns the first n meaningful lines, skipping
// frontmatter, empty lines, and HTML comments, and
// truncating long lines.
//
// The sort subpackage orders files for display.
// [sort.FilesByPriority] sorts by a configured read
// priority (CONSTITUTION first, then TASKS, etc.).
// [sort.RecentFiles] returns the n most recently
// modified files.
//
// # Data Flow
//
// The cmd/ layer loads a context entity and passes it
// to out.PersistStatusJSON or out.PersistStatusText.
// These call into sort and preview for ordering and
// content extraction. The write/status package handles
// final rendering.
package core
