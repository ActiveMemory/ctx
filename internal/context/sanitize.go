//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import "github.com/ActiveMemory/ctx/internal/config"

// effectivelyEmpty checks if a file only contains headers and whitespace.
//
// Parameters:
//   - content: File content to check
//
// Returns:
//   - bool: True if the file has no meaningful content (only headers,
//     comments, whitespace)
func effectivelyEmpty(content []byte) bool {
	// Simple heuristic: if content is shorter than the minimum
	// meaningful length, treat as empty
	if len(content) < config.MinContentLen {
		return true
	}

	// Count non-header, non-whitespace content
	lines := 0
	contentLines := 0
	openLen := len(config.CommentOpen)
	closeLen := len(config.CommentClose)
	for _, line := range splitLines(content) {
		lines++
		trimmed := trimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed[0] == '#' {
			continue
		}
		if len(trimmed) > 0 && trimmed[0] == '-' && len(trimmed) < 5 {
			continue
		}
		// Check for HTML comment markers
		if len(trimmed) >= openLen &&
			string(trimmed[:openLen]) == config.CommentOpen {
			continue
		}
		if len(trimmed) >= closeLen &&
			string(trimmed[len(trimmed)-closeLen:]) == config.CommentClose {
			continue
		}
		contentLines++
	}

	return contentLines == 0
}
