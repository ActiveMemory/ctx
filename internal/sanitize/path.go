//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	cfgSan "github.com/ActiveMemory/ctx/internal/config/sanitize"
)

// SessionID sanitizes a session identifier for safe use in file
// paths.
//
// Strips path separators, traversal sequences, and null bytes.
// Replaces remaining unsafe characters with hyphens and limits
// length to MaxSessionIDLen bytes.
//
// Parameters:
//   - s: raw session ID from MCP client
//
// Returns:
//   - string: path-safe session ID
func SessionID(s string) string {
	// Strip null bytes.
	s = strings.ReplaceAll(s, cfgSan.NullByte, "")

	// Collapse path traversal sequences.
	s = strings.ReplaceAll(s, cfgSan.DotDot, "")
	s = strings.ReplaceAll(s, cfgSan.ForwardSlash, "")
	s = strings.ReplaceAll(s, cfgSan.Backslash, "")

	// Replace remaining unsafe chars.
	s = regex.SanSessionIDUnsafe.ReplaceAllString(
		s, cfgSan.HyphenReplace,
	)

	// Remove leading/trailing hyphens.
	s = strings.Trim(s, cfgSan.HyphenReplace)

	// Limit length.
	if len(s) > cfgSan.MaxSessionIDLen {
		s = s[:cfgSan.MaxSessionIDLen]
	}

	return s
}
