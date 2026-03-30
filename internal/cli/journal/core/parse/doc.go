//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse scans journal directories and parses entries.
//
// [ScanJournalEntries] reads all markdown files in the journal
// directory, parsing frontmatter into JournalEntry structs.
// [JournalEntry] parses a single file by path.
package parse
