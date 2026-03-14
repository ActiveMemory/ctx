//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ConsolidateToolRuns collapses consecutive turns with identical body content
// into a single turn with a count. Handles both tool-call turns
// and tool-output turns.
//
// Parameters:
//   - content: Journal entry content with potential repeated turns
//
// Returns:
//   - string: Content with consecutive identical turns collapsed
func ConsolidateToolRuns(content string) string {
	lines := strings.Split(content, token.NewlineLF)
	var out []string
	i := 0

	for i < len(lines) {
		// Check if this line is a turn header
		if !regex.TurnHeader.MatchString(strings.TrimSpace(lines[i])) {
			out = append(out, lines[i])
			i++
			continue
		}

		// Extract this turn: header + body (until next header or EOF)
		header := lines[i]
		body, bodyEnd := ExtractTurnBody(lines, i+1)

		// Count consecutive turns with identical body
		count := 1
		j := bodyEnd
		for j < len(lines) {
			if !regex.TurnHeader.MatchString(strings.TrimSpace(lines[j])) {
				break
			}
			nextBody, nextBodyEnd := ExtractTurnBody(lines, j+1)

			if nextBody != body {
				break
			}
			count++
			j = nextBodyEnd
		}

		if count > 1 {
			out = append(out, header, "", body, "",
				fmt.Sprintf("(×%d)", count), "",
			)
		} else {
			// Keep original lines (preserves blank lines as-is)
			for k := i; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
		}
		i = j
	}

	return strings.Join(out, token.NewlineLF)
}
