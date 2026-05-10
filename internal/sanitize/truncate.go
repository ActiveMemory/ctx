//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import "unicode/utf8"

// truncate returns s truncated to at most maxLen bytes, with the cut
// snapped to a UTF-8 rune boundary so the result is always valid
// UTF-8. If maxLen is 0 or negative, s is returned unchanged.
//
// When the byte at position maxLen is a UTF-8 continuation byte, the
// cut backs up to the start of that rune so the partial rune is
// dropped entirely. If the entire prefix is one large rune that would
// not fit, the result is the empty string.
//
// Parameters:
//   - s: input string to truncate
//   - maxLen: maximum byte length; 0 or negative means no truncation
//
// Returns:
//   - string: s truncated to at most maxLen bytes, valid UTF-8
func truncate(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	cut := maxLen
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut]
}
