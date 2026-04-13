//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

// truncate limits a string to maxLen bytes. If truncated, no
// ellipsis is appended — the caller controls presentation.
//
// Parameters:
//   - s: input string
//   - maxLen: maximum byte length
//
// Returns:
//   - string: input capped to maxLen bytes
func truncate(s string, maxLen int) string {
	if maxLen > 0 && len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
