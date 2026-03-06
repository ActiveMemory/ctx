//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// formatSessionMatchLines formats session matches for ambiguous query output.
//
// Parameters:
//   - matches: sessions that matched the query.
//
// Returns:
//   - []string: pre-formatted lines, one per match.
func formatSessionMatchLines(matches []*parser.Session) []string {
	lines := make([]string, 0, len(matches))
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf(
			config.TplSessionMatch,
			m.Slug,
			m.ID[:config.SessionIDShortLen],
			m.StartTime.Format(config.DateTimeFormat)),
		)
	}
	return lines
}
