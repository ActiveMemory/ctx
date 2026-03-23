//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// TurnMatch holds the result of matching a turn header line.
type TurnMatch struct {
	Num  int
	Role string
	Time string
}

// MatchTurnHeader attempts to parse a turn header from a line.
//
// Parameters:
//   - line: Raw line to match (will be trimmed)
//   - mask: Whether this line is inside a pre block
//
// Returns:
//   - *TurnMatch: Parsed turn data, or nil if not a turn header
func MatchTurnHeader(line string, masked bool) *TurnMatch {
	if masked {
		return nil
	}
	m := regex.TurnHeader.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return nil
	}
	num, _ := strconv.Atoi(m[1])
	return &TurnMatch{Num: num, Role: m[2], Time: m[3]}
}

// FindTurnBoundary scans lines from startIdx to find the boundary of the
// current turn body: the last occurrence of expectedNext turn number
// that is not inside a pre block and has a timestamp >= turnTime.
//
// Parameters:
//   - lines: All lines in the document
//   - mask: Pre-block mask (true = inside pre)
//   - startIdx: Index to start scanning from
//   - turnSeq: Sorted sequence of all turn numbers in the document
//   - turnNum: Current turn number
//   - turnTime: Current turn timestamp (for ordering)
//
// Returns:
//   - int: Index of the boundary line (or len(lines) if EOF)
func FindTurnBoundary(
	lines []string, mask []bool, startIdx int,
	turnSeq []int, turnNum int, turnTime string,
) int {
	expectedNext := NextInSequence(turnSeq, turnNum)
	boundary := len(lines)
	for j := startIdx; j < len(lines); j++ {
		nm := MatchTurnHeader(lines[j], mask[j])
		if nm != nil && nm.Num == expectedNext && nm.Time >= turnTime {
			boundary = j
		}
	}
	return boundary
}

// TrimBlankLines removes leading and trailing blank lines from a slice.
//
// Parameters:
//   - lines: Input lines
//
// Returns:
//   - []string: Trimmed lines (may be empty)
func TrimBlankLines(lines []string) []string {
	start, end := 0, len(lines)-1
	for start <= end && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	for end >= start && strings.TrimSpace(lines[end]) == "" {
		end--
	}
	if start > end {
		return nil
	}
	return lines[start : end+1]
}

// ProcessTurns iterates lines, matching turn headers with the given role,
// and delegates body processing to the provided callback. Non-matching
// lines are passed through unchanged.
//
// Parameters:
//   - content: Full document content
//   - roleKey: YAML DescKey for the role to match (e.g., DescKeyLabelToolOutput)
//   - processFn: Called with (out, body, atEOF) for each matched turn;
//     returns updated out slice
//
// Returns:
//   - string: Processed content
func ProcessTurns(
	content, roleKey string,
	processFn func(out, body []string, atEOF bool) []string,
) string {
	lines := strings.Split(content, token.NewlineLF)
	mask := PreBlockMask(lines)
	turnSeq := CollectTurnNumbers(lines)
	var out []string
	i := 0

	for i < len(lines) {
		tm := MatchTurnHeader(lines[i], mask[i])
		if tm == nil || tm.Role != roleKey {
			out = append(out, lines[i])
			i++
			continue
		}

		out = append(out, lines[i])
		i++

		boundary := FindTurnBoundary(lines, mask, i, turnSeq, tm.Num, tm.Time)
		body := lines[i:boundary]
		i = boundary

		out = processFn(out, body, i >= len(lines))
	}

	return strings.Join(out, token.NewlineLF)
}
