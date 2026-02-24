//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// collapseToolOutputs wraps long Tool Output turn bodies in collapsible
// <details> blocks. Entries exported before the collapse feature was added
// have raw multi-line tool output; this pass retroactively collapses them.
//
// Parameters:
//   - content: Journal entry content
//
// Returns:
//   - string: Content with long tool outputs wrapped in <details> tags
func collapseToolOutputs(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	var out []string
	i := 0

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		matches := config.RegExTurnHeader.FindStringSubmatch(trimmed)

		// Non-header lines pass through unchanged
		if matches == nil {
			out = append(out, lines[i])
			i++
			continue
		}

		role := matches[2]
		header := lines[i]

		// Find body boundaries (mirror extractTurnBody logic)
		bodyStart := i + 1
		if bodyStart < len(lines) &&
			strings.TrimSpace(lines[bodyStart]) == "" {
			bodyStart++
		}
		bodyEnd := bodyStart
		for bodyEnd < len(lines) {
			if config.RegExTurnHeader.MatchString(
				strings.TrimSpace(lines[bodyEnd]),
			) {
				break
			}
			bodyEnd++
		}

		// Non-tool-output turns pass through unchanged
		if role != config.LabelToolOutput {
			for k := i; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
			i = bodyEnd
			continue
		}

		// Count non-blank body lines
		nonBlank := 0
		for k := bodyStart; k < bodyEnd; k++ {
			if strings.TrimSpace(lines[k]) != "" {
				nonBlank++
			}
		}

		body := strings.TrimSpace(
			strings.Join(lines[bodyStart:bodyEnd], config.NewlineLF),
		)
		alreadyWrapped := strings.HasPrefix(body, "<details>")

		if nonBlank > config.RecallDetailsThreshold && !alreadyWrapped {
			summary := fmt.Sprintf(
				config.TplRecallDetailsSummary, nonBlank,
			)
			out = append(out, header, "")
			out = append(out,
				fmt.Sprintf(config.TplRecallDetailsOpen, summary),
			)
			out = append(out, "")
			for k := bodyStart; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
			out = append(out, config.TplRecallDetailsClose, "")
		} else {
			for k := i; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
		}

		i = bodyEnd
	}

	return strings.Join(out, config.NewlineLF)
}
