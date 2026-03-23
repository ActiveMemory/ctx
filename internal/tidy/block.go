//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tidy

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
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
//   - []entity.TaskBlock: All completed top-level task blocks found
//     (outside the Completed section)
func ParseTaskBlocks(lines []string) []entity.TaskBlock {
	var blocks []entity.TaskBlock
	inCompletedSection := false
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Track if we're in the Completed section
		if strings.HasPrefix(line, desc.Text(text.DescKeyHeadingCompleted)) {
			inCompletedSection = true
			i++
			continue
		}
		if strings.HasPrefix(line, token.HeadingLevelTwoStart) && inCompletedSection {
			inCompletedSection = false
		}

		// Skip if in the Completed section or not a checked task
		match := regex.Task.FindStringSubmatch(line)
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
func RemoveBlocksFromLines(lines []string, blocks []entity.TaskBlock) []string {
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
