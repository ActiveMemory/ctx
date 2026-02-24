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

// softWrapContent wraps long lines in source journal files to ~80 characters.
// Skips only frontmatter and table rows. Wraps everything else including
// content inside code fences — journal files are reference material, not
// executable code.
//
// Parameters:
//   - content: Journal entry content with potentially long lines
//
// Returns:
//   - string: Content with long lines soft-wrapped at word boundaries
func softWrapContent(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	var out []string
	inFrontmatter := false

	for i, line := range lines {
		// Skip frontmatter
		if i == 0 && strings.TrimSpace(line) == config.Separator {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == config.Separator {
				inFrontmatter = false
			}
			continue
		}

		// Wrap long lines (skip tables)
		if len(line) > config.JournalLineWrapWidth &&
			!strings.HasPrefix(strings.TrimSpace(line), "|") {
			out = append(out, softWrap(line, config.JournalLineWrapWidth)...)
		} else {
			out = append(out, line)
		}
	}

	return strings.Join(out, config.NewlineLF)
}

// softWrap breaks a long line at word boundaries, preserving leading indent.
//
// Parameters:
//   - line: Single line to wrap
//   - width: Target line width in characters
//
// Returns:
//   - []string: Wrapped lines preserving the original indentation
func softWrap(line string, width int) []string {
	trimmed := strings.TrimLeft(line, config.Whitespace)
	indent := line[:len(line)-len(trimmed)]

	words := strings.Fields(trimmed)
	if len(words) == 0 {
		return []string{line}
	}

	var result []string
	current := indent + words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) > width && len(current) > len(indent) {
			result = append(result, current)
			current = indent + word
		} else {
			current += " " + word
		}
	}
	result = append(result, current)
	return result
}
