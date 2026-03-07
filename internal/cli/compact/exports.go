//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import "github.com/ActiveMemory/ctx/internal/cli/compact/core"

// TaskBlock is re-exported from core for backward compatibility.
//
// See core.TaskBlock for full documentation.
type TaskBlock = core.TaskBlock

// WriteArchive delegates to core.WriteArchive for backward compatibility.
//
// Parameters:
//   - prefix: File name prefix (e.g., "tasks", "decisions", "learnings")
//   - heading: Markdown heading for new archive files
//   - content: The content to archive
//
// Returns:
//   - string: Path to the written archive file
//   - error: Non-nil if directory creation or file write fails
func WriteArchive(prefix, heading, content string) (string, error) {
	return core.WriteArchive(prefix, heading, content)
}

// ParseTaskBlocks delegates to core.ParseTaskBlocks for backward compatibility.
//
// Parameters:
//   - lines: Slice of lines from the tasks file
//
// Returns:
//   - []TaskBlock: All completed top-level task blocks found
func ParseTaskBlocks(lines []string) []TaskBlock {
	return core.ParseTaskBlocks(lines)
}

// RemoveBlocksFromLines delegates to core.RemoveBlocksFromLines for backward
// compatibility.
//
// Parameters:
//   - lines: Original lines from the file
//   - blocks: Task blocks to remove (must be sorted by StartIndex)
//
// Returns:
//   - []string: New lines with blocks removed
func RemoveBlocksFromLines(lines []string, blocks []TaskBlock) []string {
	return core.RemoveBlocksFromLines(lines, blocks)
}
