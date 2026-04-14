//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

// Reflect truncates a string and strips control characters for safe
// inclusion in error or log messages.
//
// Use this for any client-supplied value that gets reflected back in
// JSON-RPC error messages (tool names, prompt names, URIs, caller
// identifiers).
//
// Parameters:
//   - s: untrusted input string
//   - maxLen: maximum allowed length (0 = no truncation)
//
// Returns:
//   - string: truncated, control-character-free string
func Reflect(s string, maxLen int) string {
	s = StripControl(s)
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}
