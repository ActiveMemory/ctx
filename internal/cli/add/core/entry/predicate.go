//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// FileTypeIsTask reports whether fileType represents a task entry.
//
// Parameters:
//   - fileType: The type string to check
//
// Returns:
//   - bool: True if fileType is a task type
func FileTypeIsTask(fileType string) bool {
	return entry.FromUserInput(fileType) == entry.Task
}

// FileTypeIsDecision reports whether fileType represents a decision entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "decision", "decisions")
//
// Returns:
//   - bool: True if fileType is a decision type
func FileTypeIsDecision(fileType string) bool {
	return entry.FromUserInput(fileType) == entry.Decision
}

// FileTypeIsLearning reports whether fileType represents a learning entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "learning", "learnings")
//
// Returns:
//   - bool: True if fileType is a learning type
func FileTypeIsLearning(fileType string) bool {
	return entry.FromUserInput(fileType) == entry.Learning
}

// NeedsSpec reports whether task content contains design-signal
// words or exceeds the length threshold configured via .ctxrc.
//
// Parameters:
//   - content: task description text
//
// Returns:
//   - bool: true if a spec nudge should fire
func NeedsSpec(content string) bool {
	if len(content) > rc.SpecNudgeMinLen() {
		return true
	}
	lower := strings.ToLower(content)
	for _, word := range rc.SpecSignalWords() {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}
