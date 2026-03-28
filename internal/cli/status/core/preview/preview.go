//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package preview

import (
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Content returns the first n non-empty, meaningful lines
// from content.
//
// Skips empty lines, YAML frontmatter delimiters, and HTML comments.
// Truncates lines longer than 60 characters.
//
// Parameters:
//   - content: The file content to extract preview from
//   - n: Maximum number of lines to return
//
// Returns:
//   - []string: Up to n meaningful lines from the content
func Content(content string, n int) []string {
	lines := strings.Split(content, token.NewlineLF)
	var preview []string

	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip YAML frontmatter
		if trimmed == token.Separator {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			continue
		}

		// Skip HTML comments
		if strings.HasPrefix(trimmed, marker.CommentOpen) {
			continue
		}

		// Truncate long lines
		if utf8.RuneCountInString(trimmed) > stats.MaxPreviewLen {
			runes := []rune(trimmed)
			truncateAt := stats.MaxPreviewLen - utf8.RuneCountInString(
				token.Ellipsis,
			)
			trimmed = string(runes[:truncateAt]) + token.Ellipsis
		}

		preview = append(preview, trimmed)
		if len(preview) >= n {
			break
		}
	}

	return preview
}
