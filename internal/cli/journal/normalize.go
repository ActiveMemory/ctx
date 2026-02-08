//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal


import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// normalizeContent sanitizes journal Markdown for static site rendering:
//   - Strips code fence markers (eliminates nesting conflicts)
//   - Strips bold markers from tool-use lines (**Glob: *.md** -> Glob: *.md)
//   - Escapes glob-like * characters outside code blocks
//
// Heavy formatting (metadata tables, proper fence reconstruction) is left to
// the ctx-journal-normalize skill which uses AI for context-aware cleanup.
//
// Parameters:
//   - content: Raw Markdown content of a journal entry
//
// Returns:
//   - string: Sanitized content ready for static site rendering
func normalizeContent(content string) string {
	// Strip fences first — eliminates all nesting conflicts
	content = stripFences(content)

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

		// Strip bold from tool-use lines
		line = config.RegExToolBold.ReplaceAllString(line, `🔧 $1`)

		// Escape glob stars
		if !strings.HasPrefix(line, "    ") {
			line = config.RegExGlobStar.ReplaceAllString(line, `\*$1`)
		}

		out = append(out, line)
	}

	return strings.Join(out, config.NewlineLF)
}