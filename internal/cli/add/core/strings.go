//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// EndsWithNewline reports whether s ends with a newline (CRLF or LF).
//
// Parameters:
//   - s: String to check
//
// Returns:
//   - bool: True if s ends with a newline
func EndsWithNewline(s string) bool {
	return strings.HasSuffix(s, token.NewlineCRLF) ||
		strings.HasSuffix(s, token.NewlineLF)
}

// Contains reports whether content contains the header and returns its index.
//
// Parameters:
//   - content: String to search in
//   - header: Substring to find
//
// Returns:
//   - bool: True if header is found
//   - int: Index of header (-1 if not found)
func Contains(content, header string) (bool, int) {
	idx := strings.Index(content, header)
	return idx != -1, idx
}

// ContainsNewLine reports whether content contains a newline and
// returns its index.
//
// Parameters:
//   - content: String to search in
//
// Returns:
//   - bool: True if a newline is found
//   - int: Index of newline (-1 if not found)
func ContainsNewLine(content string) (bool, int) {
	lineEnd := FindNewline(content)
	return lineEnd != -1, lineEnd
}

// ContainsEndComment reports whether content contains a comment close marker.
//
// Parameters:
//   - content: String to search in
//
// Returns:
//   - bool: True if comment close marker is found
//   - int: Index of marker (-1 if not found)
func ContainsEndComment(content string) (bool, int) {
	commentEnd := strings.Index(content, marker.CommentClose)
	return commentEnd != -1, commentEnd
}

// StartsWithCtxMarker reports whether s starts with a ctx marker comment.
//
// Parameters:
//   - s: String to check
//
// Returns:
//   - bool: True if s starts with CtxMarkerStart or CtxMarkerEnd
func StartsWithCtxMarker(s string) bool {
	return strings.HasPrefix(s, marker.CtxMarkerStart) ||
		strings.HasPrefix(s, marker.CtxMarkerEnd)
}
