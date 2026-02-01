//   /    Context:                     https://ctx.ist
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

// indentLevel returns the number of leading whitespace characters.
//
// Parameters:
//   - line: The line to measure
//
// Returns:
//   - int: Number of leading whitespace characters (spaces and tabs)
func indentLevel(line string) int {
	return len(line) - len(strings.TrimLeft(line, " \t"))
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
	parentIndent := indentLevel(parentLine)

	block := TaskBlock{
		Lines:        []string{parentLine},
		StartIndex:   startIdx,
		EndIndex:     startIdx + 1,
		IsCompleted:  true, // We only call this for checked tasks
		IsArchivable: true,
		DoneTime:     parseDoneTimestamp(parentLine),
	}

	// Collect all lines that are more indented than the parent
	for i := startIdx + 1; i < len(lines); i++ {
		line := lines[i]

		// Empty lines: Include if followed by more indented content
		if strings.TrimSpace(line) == "" {
			// Look ahead to see if there's more indented content
			hasMoreContent := false
			for j := i + 1; j < len(lines); j++ {
				nextLine := lines[j]
				if strings.TrimSpace(nextLine) == "" {
					continue
				}
				if indentLevel(nextLine) > parentIndent {
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
		lineIndent := indentLevel(line)
		if lineIndent <= parentIndent {
			// Same or lower indentation - end of block
			break
		}

		// This line belongs to the block
		block.Lines = append(block.Lines, line)
		block.EndIndex = i + 1

		// Check if this is an unchecked task
		nestedMatch := config.RegExTask.FindStringSubmatch(line)
		if nestedMatch != nil && task.IsPending(nestedMatch) {
			block.IsArchivable = false
		}
	}

	return block
}

// parseDoneTimestamp extracts the #done: timestamp from a task line.
//
// Parameters:
//   - line: Task line that may contain #done:YYYY-MM-DD-HHMMSS
//
// Returns:
//   - *time.Time: Parsed time, or nil if no valid timestamp is found
func parseDoneTimestamp(line string) *time.Time {
	match := config.RegExTaskDoneTimestamp.FindStringSubmatch(line)
	if match == nil || len(match) < 2 {
		return nil
	}

	// Parse YYYY-MM-DD-HHMMSS format
	t, err := time.Parse("2006-01-02-150405", match[1])
	if err != nil {
		return nil
	}
	return &t
}
