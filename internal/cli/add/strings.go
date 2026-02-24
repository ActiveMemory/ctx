//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// endsWithNewline reports whether s ends with a newline (CRLF or LF).
//
// Parameters:
//   - s: String to check
//
// Returns:
//   - bool: True if s ends with a newline
func endsWithNewline(s string) bool {
	return strings.HasSuffix(s, config.NewlineCRLF) ||
		strings.HasSuffix(s, config.NewlineLF)
}

// contains reports whether content contains the header and returns its index.
//
// Parameters:
//   - content: String to search in
//   - header: Substring to find
//
// Returns:
//   - bool: True if header is found
//   - int: Index of header (-1 if not found)
func contains(content, header string) (bool, int) {
	idx := strings.Index(content, header)
	return idx != -1, idx
}

// containsNewLine reports whether content contains a newline and
// returns its index.
//
// Parameters:
//   - content: String to search in
//
// Returns:
//   - bool: True if a newline is found
//   - int: Index of newline (-1 if not found)
func containsNewLine(content string) (bool, int) {
	lineEnd := findNewline(content)
	return lineEnd != -1, lineEnd
}

// containsEndComment reports whether content contains a comment close marker.
//
// Parameters:
//   - content: String to search in
//
// Returns:
//   - bool: True if comment close marker is found
//   - int: Index of marker (-1 if not found)
func containsEndComment(content string) (bool, int) {
	commentEnd := strings.Index(content, config.CommentClose)
	return commentEnd != -1, commentEnd
}

// startsWithCtxMarker reports whether s starts with a ctx marker comment.
//
// Parameters:
//   - s: String to check
//
// Returns:
//   - bool: True if s starts with CtxMarkerStart or CtxMarkerEnd
func startsWithCtxMarker(s string) bool {
	return strings.HasPrefix(s, config.CtxMarkerStart) ||
		strings.HasPrefix(s, config.CtxMarkerEnd)
}
