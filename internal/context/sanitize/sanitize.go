//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"bytes"

	"github.com/ActiveMemory/ctx/internal/config/content"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// EffectivelyEmpty checks if a file only contains headers and whitespace.
//
// Parameters:
//   - c: File content to check
//
// Returns:
//   - bool: True if the file has no meaningful content (only headers,
//     comments, whitespace)
func EffectivelyEmpty(c []byte) bool {
	// Simple heuristic: if the content is shorter than the minimum
	// meaningful length, treat as empty
	if len(c) < content.MinLen {
		return true
	}

	// Count non-header, non-whitespace content
	lines := 0
	contentLines := 0
	openLen := len(marker.CommentOpen)
	closeLen := len(marker.CommentClose)
	for _, line := range bytes.Split(c, []byte(token.NewlineLF)) {
		lines++
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed[0] == token.PrefixHeading[0] {
			continue
		}
		if len(trimmed) > 0 && trimmed[0] == token.Dash[0] &&
			len(trimmed) < token.MaxSeparatorLen {
			continue
		}
		// Check for HTML comment markers
		if len(trimmed) >= openLen &&
			string(trimmed[:openLen]) == marker.CommentOpen {
			continue
		}
		if len(trimmed) >= closeLen &&
			string(trimmed[len(trimmed)-closeLen:]) == marker.CommentClose {
			continue
		}
		contentLines++
	}

	return contentLines == 0
}
