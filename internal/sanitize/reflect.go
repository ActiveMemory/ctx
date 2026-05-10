//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

// Reflect strips control characters from s and truncates to maxLen
// bytes. Used when reflecting untrusted input back in error messages
// to prevent log injection.
//
// Parameters:
//   - s: untrusted input string to sanitize for reflection
//   - maxLen: maximum byte length; 0 or negative means no truncation
//
// Returns:
//   - string: s with control characters removed and length capped
func Reflect(s string, maxLen int) string {
	s = StripControl(s)
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}
