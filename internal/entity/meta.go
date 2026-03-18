//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

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
