//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import "github.com/ActiveMemory/ctx/internal/entity"

// EntryBlock represents a parsed entry block from a knowledge file
// (DECISIONS.md or LEARNINGS.md).
//
// Fields:
//   - Entry: The parsed header metadata (timestamp, date, title)
//   - Lines: All lines belonging to this entry (header + body)
//   - StartIndex: Zero-based line index where this entry starts
//   - EndIndex: Zero-based line index where this entry ends (exclusive)
type EntryBlock struct {
	Entry      entity.IndexEntry
	Lines      []string
	StartIndex int
	EndIndex   int
}
