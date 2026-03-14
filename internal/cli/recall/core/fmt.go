//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// FormatSessionMatchLines formats session matches for ambiguous query output.
//
// Parameters:
//   - matches: sessions that matched the query.
//
// Returns:
//   - []string: pre-formatted lines, one per match.
func FormatSessionMatchLines(matches []*parser.Session) []string {
	lines := make([]string, 0, len(matches))
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf(
			assets.TplSessionMatch,
			m.Slug,
			m.ID[:journal.SessionIDShortLen],
			m.StartTime.Format(time.DateTimeFormat)),
		)
	}
	return lines
}
