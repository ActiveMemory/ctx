//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

// Entry represents a parsed entry header from a context file.
//
// Fields:
//   - Timestamp: Full timestamp (YYYY-MM-DD-HHMMSS)
//   - Date: Date only (YYYY-MM-DD)
//   - Title: Entry title
type Entry struct {
	Timestamp string
	Date      string
	Title     string
}

// EntryBlock represents a parsed entry block from a knowledge file
// (DECISIONS.md or LEARNINGS.md).
//
// Fields:
//   - Entry: The parsed header metadata (timestamp, date, title)
//   - Lines: All lines belonging to this entry (header + body)
//   - StartIndex: Zero-based line index where this entry starts
//   - EndIndex: Zero-based line index where this entry ends (exclusive)
type EntryBlock struct {
	Entry      Entry
	Lines      []string
	StartIndex int
	EndIndex   int
}
