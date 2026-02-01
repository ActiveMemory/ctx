//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

// skipNewline advances pos past a newline (CRLF or LF) if present.
//
// Parameters:
//   - s: String to scan
//   - pos: Current position in s
//
// Returns:
//   - int: New position (unchanged if no newline at pos)
func skipNewline(s string, pos int) int {
	if pos >= len(s) {
		return pos
	}
	if pos+1 < len(s) && s[pos] == '\r' && s[pos+1] == '\n' {
		return pos + 2
	}
	if s[pos] == '\n' {
		return pos + 1
	}
	return pos
}

// skipWhitespace advances pos past any whitespace (space, tab, newline).
//
// Parameters:
//   - s: String to scan
//   - pos: Current position in s
//
// Returns:
//   - int: New position after skipping whitespace
func skipWhitespace(s string, pos int) int {
	for pos < len(s) {
		if n := skipNewline(s, pos); n > pos {
			pos = n
		} else if s[pos] == ' ' || s[pos] == '\t' {
			pos++
		} else {
			break
		}
	}
	return pos
}

// findNewline returns the index of the first newline (CRLF or LF) in s.
//
// Parameters:
//   - s: String to search
//
// Returns:
//   - int: Index of first newline (-1 if not found)
func findNewline(s string) int {
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && s[i] == '\r' && s[i+1] == '\n' {
			return i
		}
		if s[i] == '\n' {
			return i
		}
	}
	return -1
}
