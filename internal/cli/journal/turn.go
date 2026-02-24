//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// extractTurnBody extracts the body text from lines[start:] until the next
// turn header. Skips a leading blank line.
//
// Parameters:
//   - lines: All lines of the journal entry
//   - start: Index of the first line after the turn header
//
// Returns:
//   - string: Trimmed body content
//   - int: Index one past the last body line
func extractTurnBody(lines []string, start int) (string, int) {
	bodyStart := start
	// Skip blank line after header
	if bodyStart < len(lines) && strings.TrimSpace(lines[bodyStart]) == "" {
		bodyStart++
	}
	// Collect body until next turn header
	bodyEnd := bodyStart
	for bodyEnd < len(lines) {
		if config.RegExTurnHeader.MatchString(strings.TrimSpace(lines[bodyEnd])) {
			break
		}
		bodyEnd++
	}
	// Trim trailing blank lines for comparison
	body := strings.TrimSpace(
		strings.Join(lines[bodyStart:bodyEnd], config.NewlineLF),
	)
	return body, bodyEnd
}

// mergeConsecutiveTurns merges back-to-back turns from the same role into a
// single turn. Keeps the first header and concatenates all bodies. This reduces
// noise from sequences like 4 consecutive Assistant turns each with a single
// tool call.
//
// Parameters:
//   - content: Journal entry content with potential consecutive same-role turns
//
// Returns:
//   - string: Content with consecutive same-role turns merged
func mergeConsecutiveTurns(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	var out []string
	i := 0

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		matches := config.RegExTurnHeader.FindStringSubmatch(trimmed)
		if matches == nil {
			out = append(out, lines[i])
			i++
			continue
		}

		role := matches[2]
		header := lines[i]

		// Collect body from this and all consecutive same-role turns,
		// explicitly skipping intermediate headers.
		var body []string
		j := i + 1
		for {
			// Collect body lines until the next header or EOF
			for j < len(lines) {
				if config.RegExTurnHeader.MatchString(strings.TrimSpace(lines[j])) {
					break
				}
				body = append(body, lines[j])
				j++
			}
			// Check if the next turn has the same role
			if j >= len(lines) {
				break
			}
			nextMatches := config.RegExTurnHeader.FindStringSubmatch(
				strings.TrimSpace(lines[j]),
			)
			if nextMatches == nil || nextMatches[2] != role {
				break
			}
			// Same role — skip the header, continue collecting body
			j++
		}

		out = append(out, header)
		out = append(out, body...)
		i = j
	}

	return strings.Join(out, config.NewlineLF)
}
