//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heading

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// minHeadingLevel is the shallowest heading the projector emits. Level 1
// (`# Title`) is the file's own H1 and is intentionally excluded; entries and
// sections start at `##`.
const minHeadingLevel = 2

// Headings projects the ATX Markdown headings out of content, in file order.
//
// It is the generic counterpart to ParseHeaders: where ParseHeaders matches
// only timestamped entry headers (`## [YYYY-MM-DD-HHMMSS] Title`), Headings
// matches any `##`/`###`/… heading up to maxDepth. This is what lets one
// projector serve DECISIONS/LEARNINGS, CONVENTIONS, and TASKS (`## Phase …`)
// with the same code path.
//
// Fenced code blocks are skipped: a line like "## not a heading" inside a
// ```` ``` ```` fence is code text, not a heading, and is never projected.
//
// Parameters:
//   - content: Full file content.
//   - maxDepth: Deepest heading level to include (2 = `##` only; 3 adds
//     `###`; a value below 2 is treated as 2, the shallowest entry level).
//
// Returns:
//   - []entity.Heading: Projected headings in file order (may be empty).
func Headings(content string, maxDepth int) []entity.Heading {
	if maxDepth < minHeadingLevel {
		maxDepth = minHeadingLevel
	}

	var out []entity.Heading
	inFence := false
	inComment := false

	for _, line := range strings.Split(content, token.NewlineLF) {
		// Skip content inside a multi-line HTML comment (e.g. the
		// `<!-- DECISION FORMATS ... -->` legend), which can contain
		// example `##` headings that are documentation, not entries.
		if inComment {
			if strings.Contains(line, marker.CommentClose) {
				inComment = false
			}
			continue
		}
		// A fence marker line toggles fenced state and is never a heading.
		if regex.CodeFenceLine.MatchString(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		// A line that opens an HTML comment is not a heading. If the comment
		// also closes on the same line it is inline (skip only this line);
		// otherwise a multi-line comment block begins.
		if idx := strings.Index(line, marker.CommentOpen); idx >= 0 {
			if !strings.Contains(line[idx+len(marker.CommentOpen):], marker.CommentClose) {
				inComment = true
			}
			continue
		}

		m := regex.MarkdownHeading.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		level := len(m[1])
		if level < minHeadingLevel || level > maxDepth {
			continue
		}
		out = append(out, entity.Heading{
			Level: level,
			Text:  strings.TrimSpace(m[2]),
		})
	}

	return out
}
