//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lock manages journal entry lock state.
//
// [MatchJournalFiles] finds journal files matching a pattern.
// [MultipartBase] extracts the base name from multipart filenames.
// [UpdateFrontmatter] sets or clears the locked: field in a
// file's YAML frontmatter.
package lock
