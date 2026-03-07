//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "github.com/ActiveMemory/ctx/internal/config"

// FileTypeIsTask reports whether fileType represents a task entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "task", "tasks")
//
// Returns:
//   - bool: True if fileType is a task type
func FileTypeIsTask(fileType string) bool {
	return config.UserInputToEntry(fileType) == config.EntryTask
}

// FileTypeIsDecision reports whether fileType represents a decision entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "decision", "decisions")
//
// Returns:
//   - bool: True if fileType is a decision type
func FileTypeIsDecision(fileType string) bool {
	return config.UserInputToEntry(fileType) == config.EntryDecision
}

// FileTypeIsLearning reports whether fileType represents a learning entry.
//
// Parameters:
//   - fileType: The type string to check (e.g., "learning", "learnings")
//
// Returns:
//   - bool: True if fileType is a learning type
func FileTypeIsLearning(fileType string) bool {
	return config.UserInputToEntry(fileType) == config.EntryLearning
}
