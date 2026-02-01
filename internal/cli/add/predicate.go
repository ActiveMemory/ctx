//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import "github.com/ActiveMemory/ctx/internal/config"

// fileTypeIsTask reports whether fileType represents a task entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "task", "tasks")
//
// Returns:
//   - bool: True if fileType is a task type
func fileTypeIsTask(fileType string) bool {
	return config.UserInputToEntry(fileType) == config.EntryTask
}

// fileTypeIsDecision reports whether fileType represents a decision entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "decision", "decisions")
//
// Returns:
//   - bool: True if fileType is a decision type
func fileTypeIsDecision(fileType string) bool {
	return config.UserInputToEntry(fileType) == config.EntryDecision
}

// fileTypeIsLearning reports whether fileType represents a learning entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "learning", "learnings")
//
// Returns:
//   - bool: True if fileType is a learning type
func fileTypeIsLearning(fileType string) bool {
	return config.UserInputToEntry(fileType) == config.EntryLearning
}
