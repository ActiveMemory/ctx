//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// ExtractFrontmatter returns the YAML frontmatter block from content,
// including the delimiters and trailing newline.
//
// Parameters:
//   - content: Raw Markdown content potentially starting with frontmatter
//
// Returns:
//   - string: The frontmatter block including delimiters, or "" if not found
func ExtractFrontmatter(content string) string {
	nl := config.NewlineLF
	fmOpen := config.Separator + nl
	fmClose := nl + config.Separator + nl

	if !strings.HasPrefix(content, fmOpen) {
		return ""
	}
	end := strings.Index(content[len(fmOpen):], fmClose)
	if end < 0 {
		return ""
	}
	return content[:len(fmOpen)+end+len(fmClose)]
}

// StripFrontmatter removes the YAML frontmatter block from content,
// returning the remaining content. If no frontmatter is found, the
// original content is returned unchanged.
//
// Parameters:
//   - content: Raw Markdown content potentially starting with frontmatter
//
// Returns:
//   - string: Content without frontmatter
func StripFrontmatter(content string) string {
	fm := ExtractFrontmatter(content)
	if fm == "" {
		return content
	}
	return strings.TrimLeft(content[len(fm):], config.NewlineLF)
}
