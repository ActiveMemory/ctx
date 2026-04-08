//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package wrap

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ContextFile wraps long lines in a context file (.context/*.md).
// Handles markdown list continuation with 2-space indent. Skips
// frontmatter, headings, tables, and HTML comments.
//
// Parameters:
//   - content: Context file content with potentially long lines
//   - width: Target line width in characters
//
// Returns:
//   - string: Content with long lines soft-wrapped at word boundaries
func ContextFile(content string, width int) string {
	lines := strings.Split(content, token.NewlineLF)
	var out []string
	inFrontmatter := false
	inComment := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track frontmatter blocks.
		if i == 0 && trimmed == token.Separator {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if trimmed == token.Separator {
				inFrontmatter = false
			}
			continue
		}

		// Track HTML comment blocks.
		if strings.HasPrefix(trimmed, marker.CommentOpen) {
			inComment = true
		}
		if inComment {
			out = append(out, line)
			if strings.Contains(line, marker.CommentClose) {
				inComment = false
			}
			continue
		}

		// Skip headings, tables, and short lines.
		if len(line) <= width ||
			strings.HasPrefix(trimmed, marker.HeadingPrefix) ||
			strings.HasPrefix(trimmed, marker.TablePipe) {
			out = append(out, line)
			continue
		}

		// Wrap list items with 2-space continuation indent.
		if strings.HasPrefix(trimmed, marker.ListDash) {
			out = append(out, ListItem(line, width)...)
			continue
		}

		// Regular paragraph text.
		out = append(out, Soft(line, width)...)
	}

	return strings.Join(out, token.NewlineLF)
}

// ListItem wraps a markdown list line, using 2-space continuation
// indent for wrapped lines.
//
// Parameters:
//   - line: Single list-item line (e.g., "- [ ] long text...")
//   - width: Target line width in characters
//
// Returns:
//   - []string: First line with original prefix, continuation lines
//     indented with 2 spaces
func ListItem(line string, width int) []string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{line}
	}

	var result []string
	current := words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) > width && len(current) > 0 {
			result = append(result, current)
			current = marker.ListContinuationIndent + word
		} else {
			current += token.Space + word
		}
	}
	result = append(result, current)
	return result
}
