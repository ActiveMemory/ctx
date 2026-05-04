//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

// truncate returns s truncated to maxLen bytes.
// If maxLen is 0 or negative, s is returned unchanged.
//
// Parameters:
//   - s: input string to truncate
//   - maxLen: maximum byte length; 0 or negative means no truncation
//
// Returns:
//   - string: s truncated to at most maxLen bytes
func truncate(s string, maxLen int) string {
	if maxLen > 0 && len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
