//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// contentPreview returns the first n non-empty, meaningful lines
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
func contentPreview(content string, n int) []string {
	lines := strings.Split(content, config.NewlineLF)
	var preview []string

	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip YAML frontmatter
		if trimmed == config.Separator {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			continue
		}

		// Skip HTML comments
		if strings.HasPrefix(trimmed, config.CommentOpen) {
			continue
		}

		// Truncate long lines
		if len(trimmed) > config.MaxPreviewLen {
			truncateAt := config.MaxPreviewLen - len(config.Ellipsis)
			trimmed = trimmed[:truncateAt] + config.Ellipsis
		}

		preview = append(preview, trimmed)
		if len(preview) >= n {
			break
		}
	}

	return preview
}
