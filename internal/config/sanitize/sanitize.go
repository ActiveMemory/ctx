//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

// Sanitize-layer string and length constants.
const (
	// NullByte is the null character stripped from untrusted input.
	NullByte = "\x00"

	// DotDot is a path traversal sequence.
	DotDot = ".."

	// ForwardSlash is the forward slash stripped from session IDs.
	ForwardSlash = "/"

	// Backslash is the backslash stripped from session IDs.
	Backslash = "\\"

	// HyphenReplace is the replacement character for unsafe
	// session ID characters.
	HyphenReplace = "-"

	// EscapePrefix is the backslash prefix for escaping Markdown
	// structural patterns.
	EscapePrefix = `\`

	// MaxSessionIDLen is the maximum byte length for a session
	// identifier.
	MaxSessionIDLen = 128
)
