//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

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

// ParseEntryBlocks splits file content into discrete entry blocks.
//
// Each block starts at a "## [YYYY-MM-DD-HHMMSS] Title" header and extends
// to the line before the next entry header or end of content.
//
// Parameters:
//   - content: The full file content
//
// Returns:
//   - []EntryBlock: Parsed entry blocks in file order (may be empty)
func ParseEntryBlocks(content string) []EntryBlock {
	if content == "" {
		return nil
	}

	lines := strings.Split(content, config.NewlineLF)
	var blocks []EntryBlock

	// Find all entry header positions
	type headerPos struct {
		lineIdx int
		entry   Entry
	}
	var headers []headerPos

	for i, line := range lines {
		matches := config.RegExEntryHeader.FindStringSubmatch(line)
		if len(matches) == 4 {
			headers = append(headers, headerPos{
				lineIdx: i,
				entry: Entry{
					Timestamp: matches[1] + "-" + matches[2],
					Date:      matches[1],
					Title:     matches[3],
				},
			})
		}
	}

	if len(headers) == 0 {
		return nil
	}

	for i, h := range headers {
		var endIdx int
		if i+1 < len(headers) {
			endIdx = headers[i+1].lineIdx
		} else {
			endIdx = len(lines)
		}

		// Trim trailing blank lines from the block
		for endIdx > h.lineIdx+1 && strings.TrimSpace(lines[endIdx-1]) == "" {
			endIdx--
		}

		blocks = append(blocks, EntryBlock{
			Entry:      h.entry,
			Lines:      lines[h.lineIdx:endIdx],
			StartIndex: h.lineIdx,
			EndIndex:   endIdx,
		})
	}

	return blocks
}

// IsSuperseded checks whether this entry has been marked as superseded.
//
// An entry is superseded when its body contains a line starting with
// "~~Superseded" (strikethrough prefix).
//
// Returns:
//   - bool: True if the entry contains a superseded marker
func (eb *EntryBlock) IsSuperseded() bool {
	for _, line := range eb.Lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~Superseded") {
			return true
		}
	}
	return false
}

// OlderThan checks whether this entry's timestamp is older than the given
// number of days from now.
//
// Parameters:
//   - days: Number of days threshold
//
// Returns:
//   - bool: True if the entry is older than the threshold
func (eb *EntryBlock) OlderThan(days int) bool {
	// Parse the date portion of the timestamp (YYYY-MM-DD)
	entryDate, err := time.ParseInLocation("2006-01-02", eb.Entry.Date, time.Local)
	if err != nil {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return entryDate.Before(threshold)
}

// BlockContent joins the entry's lines into a single string.
//
// Returns:
//   - string: The full entry content with lines joined by newlines
func (eb *EntryBlock) BlockContent() string {
	return strings.Join(eb.Lines, config.NewlineLF)
}

// RemoveEntryBlocks removes the specified entry blocks from file content
// and cleans up excess blank lines.
//
// Parameters:
//   - content: The full file content
//   - blocks: Entry blocks to remove (must have valid StartIndex/EndIndex)
//
// Returns:
//   - string: Content with the specified blocks removed
func RemoveEntryBlocks(content string, blocks []EntryBlock) string {
	if len(blocks) == 0 {
		return content
	}

	lines := strings.Split(content, config.NewlineLF)

	// Build a set of line indices to remove
	removeSet := make(map[int]bool)
	for _, b := range blocks {
		for i := b.StartIndex; i < b.EndIndex; i++ {
			removeSet[i] = true
		}
	}

	// Build new lines, skipping removed indices
	var newLines []string
	for i, line := range lines {
		if !removeSet[i] {
			newLines = append(newLines, line)
		}
	}

	// Clean up excess blank lines (collapse 3+ consecutive blanks to 2)
	var cleaned []string
	blankCount := 0
	for _, line := range newLines {
		if strings.TrimSpace(line) == "" {
			blankCount++
			if blankCount <= 2 {
				cleaned = append(cleaned, line)
			}
		} else {
			blankCount = 0
			cleaned = append(cleaned, line)
		}
	}

	// Trim trailing blank lines
	for len(cleaned) > 0 && strings.TrimSpace(cleaned[len(cleaned)-1]) == "" {
		cleaned = cleaned[:len(cleaned)-1]
	}

	if len(cleaned) == 0 {
		return ""
	}

	return strings.Join(cleaned, config.NewlineLF) + config.NewlineLF
}
