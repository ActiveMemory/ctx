//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status provides terminal output for the
// context status command (ctx status).
//
// The status command displays a summary of the
// context directory: its path, file inventory,
// token estimates, and recent activity.
//
// # Header
//
// [Header] renders the context directory path with
// the total file count and estimated token count.
// Token counts are formatted with thousand
// separators for readability.
//
// # File Listing
//
// [FileItem] renders a single context file entry.
// In compact mode it shows a status indicator,
// file name, and status text. In verbose mode it
// adds token count, byte size, and a content
// preview of the first few lines.
//
// # Activity
//
// [Activity] renders the recent session activity
// section with a header and a list of entries,
// each showing a file name and a human-readable
// time-ago string.
//
// # Data Types
//
// [FileInfo] carries pre-computed display data for
// a single file: indicator glyph, name, status
// text, token count, byte size, and preview lines.
// [ActivityInfo] carries a file name and its
// time-ago string. Both types keep business logic
// out of the output functions.
package status
