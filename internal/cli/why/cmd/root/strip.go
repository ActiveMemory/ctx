//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"regexp"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// linkRe matches Markdown links with relative .md targets.
var linkRe = regexp.MustCompile(`\[([^\]]+)\]\([^\)]*\.md[^\)]*\)`)

// imageRe matches Markdown image lines.
var imageRe = regexp.MustCompile(`^\s*!\[.*\]\(.*\)\s*$`)

// StripMkDocs removes MkDocs-specific syntax from Markdown content so it
// reads cleanly in the terminal.
//
// Handles:
//   - YAML frontmatter (--- blocks) -- removed
//   - Image refs (![alt](path)) -- line removed
//   - Admonitions (!!! type "Title") -- converted to blockquote with bold title
//   - Tab markers (=== "Name") -- converted to bold name; body dedented
//   - Relative .md links ([text](file.md)) -- kept as text only
//
// Parameters:
//   - content: Raw Markdown with MkDocs syntax
//
// Returns:
//   - string: Cleaned Markdown suitable for terminal display
func StripMkDocs(content string) string {
	lines := strings.Split(content, config.NewlineLF)
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
			title := ExtractAdmonitionTitle(line)
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
			title := ExtractTabTitle(line)
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

	return strings.Join(result, config.NewlineLF)
}

// ExtractAdmonitionTitle pulls the quoted title from an admonition line.
// e.g., `!!! note "Title"` -> "Title"
//
// Parameters:
//   - line: The admonition line to parse
//
// Returns:
//   - string: The extracted title, or empty string if no valid title found
func ExtractAdmonitionTitle(line string) string {
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

// ExtractTabTitle pulls the quoted title from a tab marker line.
// e.g., `=== "Name"` -> "Name"
//
// Parameters:
//   - line: The tab marker line to parse
//
// Returns:
//   - string: The extracted title, or empty string if no valid title found
func ExtractTabTitle(line string) string {
	return ExtractAdmonitionTitle(line)
}
