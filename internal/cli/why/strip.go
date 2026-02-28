//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"regexp"
	"strings"
)

// linkRe matches Markdown links with relative .md targets.
var linkRe = regexp.MustCompile(`\[([^\]]+)\]\([^\)]*\.md[^\)]*\)`)

// imageRe matches Markdown image lines.
var imageRe = regexp.MustCompile(`^\s*!\[.*\]\(.*\)\s*$`)

// StripMkDocs removes MkDocs-specific syntax from Markdown content so it
// reads cleanly in the terminal.
//
// Handles:
//   - YAML frontmatter (--- blocks) — removed
//   - Image refs (![alt](path)) — line removed
//   - Admonitions (!!! type "Title") — converted to blockquote with bold title
//   - Tab markers (=== "Name") — converted to bold name; body dedented
//   - Relative .md links ([text](file.md)) — kept as text only
//
// Parameters:
//   - content: Raw Markdown with MkDocs syntax
//
// Returns:
//   - string: Cleaned Markdown suitable for terminal display
func StripMkDocs(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	// Strip YAML frontmatter.
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				lines = lines[i+1:]
				break
			}
		}
	}

	inAdmonition := false
	inTab := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Skip image lines.
		if imageRe.MatchString(line) {
			continue
		}

		// Admonition start: !!! type "Title"
		if strings.HasPrefix(strings.TrimSpace(line), "!!!") {
			inAdmonition = true
			title := extractAdmonitionTitle(line)
			if title != "" {
				result = append(result, "> **"+title+"**")
			}
			continue
		}

		// Inside admonition: dedent 4-space body.
		if inAdmonition {
			if strings.HasPrefix(line, "    ") {
				result = append(result, "> "+line[4:])
				continue
			}
			// End of admonition body.
			inAdmonition = false
		}

		// Tab marker: === "Name"
		if strings.HasPrefix(strings.TrimSpace(line), "=== ") {
			inTab = true
			title := extractTabTitle(line)
			if title != "" {
				result = append(result, "**"+title+"**")
			}
			continue
		}

		// Inside tab: dedent 4-space body.
		if inTab {
			if strings.HasPrefix(line, "    ") {
				result = append(result, line[4:])
				continue
			}
			if strings.TrimSpace(line) == "" {
				result = append(result, "")
				continue
			}
			// Non-indented, non-empty line ends the tab block.
			inTab = false
		}

		// Strip relative .md links, keep display text.
		line = linkRe.ReplaceAllString(line, "$1")

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// extractAdmonitionTitle pulls the quoted title from an admonition line.
// e.g., `!!! note "Title"` → "Title"
func extractAdmonitionTitle(line string) string {
	idx := strings.Index(line, `"`)
	if idx < 0 {
		return ""
	}
	end := strings.LastIndex(line, `"`)
	if end <= idx {
		return ""
	}
	return line[idx+1 : end]
}

// extractTabTitle pulls the quoted title from a tab marker line.
// e.g., `=== "Name"` → "Name"
func extractTabTitle(line string) string {
	return extractAdmonitionTitle(line)
}
