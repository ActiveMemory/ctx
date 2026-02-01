//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.

package context

import "time"

// FileInfo represents metadata about a context file.
//
// Fields:
//   - Name: Filename (e.g., "TASKS.md")
//   - Path: Full path to the file
//   - Size: File size in bytes
//   - ModTime: Last modification time
//   - Content: Raw file content
//   - IsEmpty: True if the file has no meaningful content
//     (only headers/whitespace)
//   - Tokens: Estimated token count for the content
//   - Summary: Brief description generated from the content
type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
	Content []byte
	IsEmpty bool
	Tokens  int
	Summary string
}

// Context represents the loaded context from a .context/ directory.
//
// Fields:
//   - Dir: Path to the context directory
//   - Files: All loaded context files with their metadata
//   - TotalTokens: Sum of estimated tokens across all files
//   - TotalSize: Sum of file sizes in bytes
type Context struct {
	Dir         string
	Files       []FileInfo
	TotalTokens int
	TotalSize   int64
}

// NotFoundError is returned when the context directory doesn't exist.
type NotFoundError struct {
	Dir string
}

// Error implements the error interface for NotFoundError.
//
// Returns:
//   - string: Error message including the missing directory path
func (e *NotFoundError) Error() string {
	return "context directory not found: " + e.Dir
}
