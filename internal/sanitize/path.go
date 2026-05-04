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

// SessionID converts an arbitrary string into a path-safe session
// identifier by stripping dangerous characters and truncating to
// [cfgSan.MaxSessionIDLen].
//
// Parameters:
//   - s: raw session identifier from untrusted input
//
// Returns:
//   - string: safe session identifier suitable for use in file paths
func SessionID(s string) string {
	s = strings.ReplaceAll(s, cfgSan.NullByte, "")
	s = strings.ReplaceAll(s, cfgSan.DotDot, "")
	s = strings.ReplaceAll(s, cfgSan.ForwardSlash, "")
	s = strings.ReplaceAll(s, cfgSan.Backslash, "")
	s = regex.SanSessionIDUnsafe.ReplaceAllString(s, cfgSan.HyphenReplace)
	s = strings.Trim(s, cfgSan.HyphenReplace)
	if len(s) > cfgSan.MaxSessionIDLen {
		s = s[:cfgSan.MaxSessionIDLen]
	}
	return s
}
