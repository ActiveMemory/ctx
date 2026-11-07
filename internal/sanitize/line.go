//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"strings"
	"unicode"

	cfgSan "github.com/ActiveMemory/ctx/internal/config/sanitize"
)

// Line flattens s to a single line: it removes every control
// character, including the tab, line feed, and carriage return that
// [StripControl] keeps, and the Unicode line/paragraph separators,
// then truncates to maxLen bytes on a UTF-8 rune boundary. Used for
// single-line values written into a Markdown line or heading, where a
// line break would start a new one.
//
// Parameters:
//   - s: untrusted single-line input
//   - maxLen: maximum byte length; 0 or negative means no truncation
//
// Returns:
//   - string: s on one line, control characters removed and length
//     capped
func Line(s string, maxLen int) string {
	return truncate(strings.Map(func(r rune) rune {
		if unicode.IsControl(r) ||
			r == cfgSan.LineSeparator || r == cfgSan.ParagraphSeparator {
			return -1
		}
		return r
	}, s), maxLen)
}
