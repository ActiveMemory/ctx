//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// CollapseToolOutputs wraps long Tool Output turn bodies in collapsible
// <details> blocks. Entries exported before the collapse feature was added
// have raw multi-line tool output; this pass retroactively collapses them.
//
// Parameters:
//   - content: Journal entry content
//
// Returns:
//   - string: Content with long tool outputs wrapped in <details> tags
func CollapseToolOutputs(content string) string {
	lines := strings.Split(content, token.NewlineLF)
	var out []string
	i := 0

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		matches := regex.TurnHeader.FindStringSubmatch(trimmed)

		// Non-header lines pass through unchanged
		if matches == nil {
			out = append(out, lines[i])
			i++
			continue
		}

		role := matches[2]
		header := lines[i]

		// Find body boundaries (mirror ExtractTurnBody logic)
		bodyStart := i + 1
		if bodyStart < len(lines) &&
			strings.TrimSpace(lines[bodyStart]) == "" {
			bodyStart++
		}
		bodyEnd := bodyStart
		for bodyEnd < len(lines) {
			if regex.TurnHeader.MatchString(
				strings.TrimSpace(lines[bodyEnd]),
			) {
				break
			}
			bodyEnd++
		}

		// Non-tool-output turns pass through unchanged
		if role != assets.ToolOutput {
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
			strings.Join(lines[bodyStart:bodyEnd], token.NewlineLF),
		)
		alreadyWrapped := strings.HasPrefix(body, "<details>")

		if nonBlank > journal.DetailsThreshold && !alreadyWrapped {
			summary := fmt.Sprintf(
				assets.TplRecallDetailsSummary, nonBlank,
			)
			out = append(out, header, "")
			out = append(out,
				fmt.Sprintf(assets.TplRecallDetailsOpen, summary),
			)
			out = append(out, "")
			for k := bodyStart; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
			out = append(out, assets.TplRecallDetailsClose, "")
		} else {
			for k := i; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
		}

		i = bodyEnd
	}

	return strings.Join(out, token.NewlineLF)
}
