//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"regexp"
	"strings"
)

// TaskBlock represents a task and its nested content.
//
// Fields:
//   - Lines: All lines in the block (parent + children)
//   - StartIndex: Index of first line in original content
//   - EndIndex: Index of last line (exclusive)
//   - IsCompleted: Parent task is checked
//   - IsArchivable: Completed and no unchecked children
type TaskBlock struct {
	Lines        []string
	StartIndex   int
	EndIndex     int
	IsCompleted  bool
	IsArchivable bool
}

// Patterns for task detection
var (
	// Matches checked task: "- [x] content" or "  - [x] content"
	checkedTaskPattern = regexp.MustCompile(`^(\s*)-\s*\[x]\s*(.+)$`)
	// Matches unchecked task: "- [ ] content" or "  - [ ] content"
	uncheckedTaskPattern = regexp.MustCompile(`^(\s*)-\s*\[\s*]\s*(.+)$`)
	// Matches any task (checked or unchecked)
	anyTaskPattern = regexp.MustCompile(`^(\s*)-\s*\[[x ]*]\s*(.+)$`)
)

// UncheckedTaskPattern returns the regex for matching unchecked tasks.
//
// Returns:
//   - *regexp.Regexp: Pattern matching "- [ ] content" with optional indentation
func UncheckedTaskPattern() *regexp.Regexp {
	return uncheckedTaskPattern
}

// GetIndentLevel returns the number of leading whitespace characters.
//
// Parameters:
//   - line: The line to measure
//
// Returns:
//   - int: Number of leading whitespace characters (spaces and tabs)
func GetIndentLevel(line string) int {
	return getIndentLevel(line)
}

// getIndentLevel returns the number of leading whitespace characters.
//
// Parameters:
//   - line: The line to measure
//
// Returns:
//   - int: Number of leading whitespace characters (spaces and tabs)
func getIndentLevel(line string) int {
	return len(line) - len(strings.TrimLeft(line, " \t"))
}

// ParseTaskBlocks parses content into task blocks, identifying completed tasks
// with their nested content.
//
// A task block consists of:
// - A parent task line at the top level (no indentation, e.g., "- [x] Task title")
// - All following lines that are more indented than the parent
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
//   - []TaskBlock: All completed top-level task blocks found (outside Completed section)
func ParseTaskBlocks(lines []string) []TaskBlock {
	var blocks []TaskBlock
	inCompletedSection := false
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Track if we're in the Completed section
		if strings.HasPrefix(line, "## Completed") {
			inCompletedSection = true
			i++
			continue
		}
		if strings.HasPrefix(line, "## ") && inCompletedSection {
			inCompletedSection = false
		}

		// Skip if in Completed section or not a checked task
		if inCompletedSection || !checkedTaskPattern.MatchString(line) {
			i++
			continue
		}

		// Only consider top-level tasks (no indentation) for archiving
		// Nested subtasks are part of their parent block, not independent
		if getIndentLevel(line) > 0 {
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

// parseBlockAt parses a task block starting at the given index.
//
// Parameters:
//   - lines: All lines from the file
//   - startIdx: Index of the parent task line
//
// Returns:
//   - TaskBlock: Parsed block with all nested content
func parseBlockAt(lines []string, startIdx int) TaskBlock {
	parentLine := lines[startIdx]
	parentIndent := getIndentLevel(parentLine)

	block := TaskBlock{
		Lines:       []string{parentLine},
		StartIndex:  startIdx,
		EndIndex:    startIdx + 1,
		IsCompleted: true, // We only call this for checked tasks
		IsArchivable: true,
	}

	// Collect all lines that are more indented than the parent
	for i := startIdx + 1; i < len(lines); i++ {
		line := lines[i]

		// Empty lines: include if followed by more indented content
		if strings.TrimSpace(line) == "" {
			// Look ahead to see if there's more indented content
			hasMoreContent := false
			for j := i + 1; j < len(lines); j++ {
				nextLine := lines[j]
				if strings.TrimSpace(nextLine) == "" {
					continue
				}
				if getIndentLevel(nextLine) > parentIndent {
					hasMoreContent = true
				}
				break
			}
			if hasMoreContent {
				block.Lines = append(block.Lines, line)
				block.EndIndex = i + 1
				continue
			}
			// No more indented content, stop here
			break
		}

		// Check indentation
		lineIndent := getIndentLevel(line)
		if lineIndent <= parentIndent {
			// Same or lower indentation - end of block
			break
		}

		// This line belongs to the block
		block.Lines = append(block.Lines, line)
		block.EndIndex = i + 1

		// Check if this is an unchecked task
		if uncheckedTaskPattern.MatchString(line) {
			block.IsArchivable = false
		}
	}

	return block
}

// BlockContent returns the full content of a block as a single string.
//
// Returns:
//   - string: All lines joined with newlines
func (b *TaskBlock) BlockContent() string {
	return strings.Join(b.Lines, "\n")
}

// ParentTaskText extracts just the task text from the parent line.
//
// Returns:
//   - string: Task text without checkbox prefix, empty if no lines
func (b *TaskBlock) ParentTaskText() string {
	if len(b.Lines) == 0 {
		return ""
	}
	matches := checkedTaskPattern.FindStringSubmatch(b.Lines[0])
	if len(matches) > 2 {
		return matches[2]
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
		if blockIdx < len(blocks) && i >= blocks[blockIdx].StartIndex && i < blocks[blockIdx].EndIndex {
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
