//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "github.com/ActiveMemory/ctx/internal/config"

// SkipNewline advances pos past a newline (CRLF or LF) if present.
//
// Parameters:
//   - s: String to scan
//   - pos: Current position in s
//
// Returns:
//   - int: New position (unchanged if no newline at pos)
func SkipNewline(s string, pos int) int {
	if pos >= len(s) {
		return pos
	}
	if pos+len(config.NewlineCRLF) <= len(s) &&
		s[pos] == config.NewlineCRLF[0] && s[pos+1] == config.NewlineCRLF[1] {
		return pos + len(config.NewlineCRLF)
	}
	if s[pos] == config.NewlineLF[0] {
		return pos + len(config.NewlineLF)
	}
	return pos
}

// SkipWhitespace advances pos past any whitespace (space, tab, newline).
//
// Parameters:
//   - s: String to scan
//   - pos: Current position in s
//
// Returns:
//   - int: New position after skipping whitespace
func SkipWhitespace(s string, pos int) int {
	for pos < len(s) {
		if n := SkipNewline(s, pos); n > pos {
			pos = n
		} else if s[pos] == config.Space[0] || s[pos] == config.Tab[0] {
			pos++
		} else {
			break
		}
	}
	return pos
}

// FindNewline returns the index of the first newline (CRLF or LF) in s.
//
// Parameters:
//   - s: String to search
//
// Returns:
//   - int: Index of the first newline (-1 if not found)
func FindNewline(s string) int {
	for i := 0; i < len(s); i++ {
		if i+len(config.NewlineCRLF) <= len(s) &&
			s[i] == config.NewlineCRLF[0] && s[i+1] == config.NewlineCRLF[1] {
			return i
		}
		if s[i] == config.NewlineLF[0] {
			return i
		}
	}
	return -1
}
