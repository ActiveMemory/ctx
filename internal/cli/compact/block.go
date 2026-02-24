//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/task"
)

// ParseTaskBlocks parses content into task blocks, identifying completed tasks
// with their nested content.
//
// A task block consists of:
//   - A parent task line at the top level (no indentation, e.g.,
//     "- [x] Task title")
//   - All following lines that are more indented than the parent
//
// Only top-level tasks (indent=0) are considered for archiving. Nested subtasks
// are part of their parent block and are not collected as independent blocks.
//
// A block is archivable if:
// - The parent task is checked [x]
// - No nested lines contain unchecked tasks [ ]
//
// Parameters:
//   - lines: Slice of lines from the tasks file
//
// Returns:
//   - []TaskBlock: All completed top-level task blocks found
//     (outside the Completed section)
func ParseTaskBlocks(lines []string) []TaskBlock {
	var blocks []TaskBlock
	inCompletedSection := false
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Track if we're in the Completed section
		if strings.HasPrefix(line, config.HeadingCompleted) {
			inCompletedSection = true
			i++
			continue
		}
		if strings.HasPrefix(line, config.HeadingLevelTwoStart) && inCompletedSection {
			inCompletedSection = false
		}

		// Skip if in the Completed section or not a checked task
		match := config.RegExTask.FindStringSubmatch(line)
		if inCompletedSection || match == nil || !task.Completed(match) {
			i++
			continue
		}

		// Only consider top-level tasks (no indentation) for archiving
		// Nested subtasks are part of their parent block, not independent
		if indentLevel(line) > 0 {
			i++
			continue
		}

		// Found a completed top-level task - parse the block
		block := parseBlockAt(lines, i)
		blocks = append(blocks, block)

		// Skip past this block
		i = block.EndIndex
	}

	return blocks
}

// BlockContent returns the full content of a block as a single string.
//
// Returns:
//   - string: All lines joined with newlines
func (b *TaskBlock) BlockContent() string {
	return strings.Join(b.Lines, config.NewlineLF)
}

// ParentTaskText extracts just the task text from the parent line.
//
// Returns:
//   - string: Task text without the checkbox prefix, empty if no lines
func (b *TaskBlock) ParentTaskText() string {
	if len(b.Lines) == 0 {
		return ""
	}
	match := config.RegExTask.FindStringSubmatch(b.Lines[0])
	if match != nil {
		return task.Content(match)
	}
	return ""
}

// RemoveBlocksFromLines removes the specified blocks from the lines slice.
//
// Blocks must be sorted by StartIndex in ascending order.
//
// Parameters:
//   - lines: Original lines from the file
//   - blocks: Task blocks to remove (must be sorted by StartIndex)
//
// Returns:
//   - []string: New lines with blocks removed
func RemoveBlocksFromLines(lines []string, blocks []TaskBlock) []string {
	if len(blocks) == 0 {
		return lines
	}

	var result []string
	blockIdx := 0

	for i := 0; i < len(lines); i++ {
		// Check if this line is part of a block to remove
		if blockIdx < len(blocks) &&
			i >= blocks[blockIdx].StartIndex && i < blocks[blockIdx].EndIndex {
			// Skip this line
			if i == blocks[blockIdx].EndIndex-1 {
				blockIdx++
			}
			continue
		}
		result = append(result, lines[i])
	}

	return result
}

// OlderThan checks if the task was completed more than the specified days ago.
//
// Parameters:
//   - days: Number of days threshold
//
// Returns:
//   - bool: True if DoneTime is set and older than days ago, false otherwise
func (b *TaskBlock) OlderThan(days int) bool {
	if b.DoneTime == nil {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return b.DoneTime.Before(threshold)
}
