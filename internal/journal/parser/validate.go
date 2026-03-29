//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"strings"

	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// sessionHeader checks if a line is a session header.
// (i.e. "Session: YYYY-MM-DD - Topic")
//
// Parameters:
//   - line: Trimmed line to check
//
// Returns:
//   - bool: True if the line matches a session header pattern
func sessionHeader(line string) bool {
	if !strings.HasPrefix(line, token.HeadingLevelOneStart) {
		return false
	}

	rest := line[len(token.HeadingLevelOneStart):]

	// Check for configured session prefixes (e.g., "Session:")
	for _, prefix := range rc.SessionPrefixes() {
		if strings.HasPrefix(rest, prefix) {
			return true
		}
	}

	// Check for a direct date pattern (YYYY-MM-DD)
	dash := token.Dash[0]
	if len(rest) >= cfgTime.DateMinLen && rest[cfgTime.DateHyphenPos1] == dash &&
		rest[cfgTime.DateHyphenPos2] == dash {
		return true
	}

	return false
}
