//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"regexp"
	"strconv"
	"strings"

	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/parse"
)

// ParseTurnRange parses a "START:END" or "N" turn selector into 1-based
// inclusive bounds clamped against totalTurns.
//
// Parameters:
//   - spec: Turn range string ("START:END" or "N")
//   - totalTurns: Total message count in the session (0 skips upper clamping)
//
// Returns:
//   - int: 1-based inclusive start turn
//   - int: 1-based inclusive end turn
//   - error: Non-nil if spec syntax or bounds are invalid
func ParseTurnRange(spec string, totalTurns int) (int, int, error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return 0, 0, errSession.InvalidTurnRange(spec)
	}
	parts := strings.Split(trimmed, token.Colon)
	if len(parts) > 2 {
		return 0, 0, errSession.InvalidTurnRange(spec)
	}
	start, startErr := strconv.Atoi(strings.TrimSpace(parts[0]))
	if startErr != nil || start < 1 {
		return 0, 0, errSession.InvalidTurnRange(spec)
	}
	end := start
	if len(parts) == 2 {
		var endErr error
		end, endErr = strconv.Atoi(strings.TrimSpace(parts[1]))
		if endErr != nil || end < start {
			return 0, 0, errSession.InvalidTurnRange(spec)
		}
	}
	if totalTurns > 0 && start <= totalTurns && end > totalTurns {
		end = totalTurns
	}
	return start, end, nil
}

// MatchTurns filters session messages within [start, end] against a
// compiled regular expression.
//
// Parameters:
//   - msgs: Session messages slice
//   - start: 1-based inclusive start turn
//   - end: 1-based inclusive end turn
//   - re: Compiled regular expression to match
//
// Returns:
//   - []TurnMatch: Matching turns in ascending turn order
func MatchTurns(
	msgs []entity.Message, start, end int, re *regexp.Regexp,
) []TurnMatch {
	if re == nil || len(msgs) == 0 || start < 1 || end < start {
		return nil
	}
	if start > len(msgs) {
		return nil
	}
	limit := min(end, len(msgs))
	var matches []TurnMatch
	for idx := start - 1; idx < limit; idx++ {
		msg := msgs[idx]
		var haystacks []string
		if msg.Text != "" {
			haystacks = append(haystacks, msg.Text)
		}
		for _, tu := range msg.ToolUses {
			haystacks = append(haystacks, srcFmt.ToolUse(tu))
		}
		for _, tr := range msg.ToolResults {
			if tr.Content != "" {
				haystacks = append(
					haystacks, parse.StripLineNumbers(tr.Content),
				)
			}
		}
		for _, h := range haystacks {
			loc := re.FindStringIndex(h)
			if loc != nil {
				matches = append(matches, TurnMatch{
					Turn:      idx + 1,
					Role:      turnRoleLabel(msg),
					Timestamp: msg.Timestamp.Format(time.Format),
					Snippet: extractMatchSnippet(
						h, loc, journal.GrepSnippetMaxChars,
					),
				})
				break
			}
		}
	}
	return matches
}
