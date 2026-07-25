//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/normalize"
	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/inspect"
)

// conventionTarget resolves the H2 heading a new bullet belongs under.
//
// Parameters:
//   - content: Existing file content
//   - section: Caller-supplied section name, possibly empty
//
// Returns:
//   - string: the normalized H2 heading line to target
func conventionTarget(content, section string) string {
	if section != "" {
		return normalize.ConventionSection(section)
	}
	// Unreachable from the CLI: requireSection refuses an empty --section
	// for conventions before the add-path runs. Direct callers of
	// AppendEntry fall back to the first section rather than inventing a
	// catch-all, so a bullet still lands somewhere enumerable.
	if first, ok := firstConventionSection(content); ok {
		return first
	}
	return normalize.ConventionSection(cfgDisc.SectionUnfiled)
}

// isConventionSection reports whether a line opens a staged convention
// section: an H2 heading that is not the structural "## Themes".
//
// Parameters:
//   - line: the line to test
//
// Returns:
//   - bool: true when the line opens a section
func isConventionSection(line string) bool {
	return strings.HasPrefix(line, token.HeadingLevelTwoStart) &&
		strings.TrimSpace(line) != cfgDisc.HeadingThemes
}

// firstConventionSection returns the root's first section heading,
// skipping headings inside HTML comments (authoring examples).
//
// Parameters:
//   - content: Existing file content
//
// Returns:
//   - string: the heading line
//   - bool: false when the root holds no section
func firstConventionSection(content string) (string, bool) {
	off := 0
	for _, line := range strings.Split(content, token.NewlineLF) {
		if isConventionSection(line) &&
			!ExistsInsideHTMLComment(content, off) {
			return strings.TrimRight(line, token.Whitespace), true
		}
		off += len(line) + len(token.NewlineLF)
	}
	return "", false
}

// conventionBodyStart returns the byte offset where a section's body
// begins — after its heading line and the blank line that follows — so an
// inserted bullet lands newest-first rather than jammed against the
// heading.
//
// Parameters:
//   - content: Existing file content
//   - header: the heading line to find
//
// Returns:
//   - int: offset where the section body starts
//   - bool: false when the heading is absent
func conventionBodyStart(content, header string) (int, bool) {
	off := 0
	for _, line := range strings.Split(content, token.NewlineLF) {
		next := off + len(line) + len(token.NewlineLF)
		if strings.TrimRight(line, token.Whitespace) == header &&
			!ExistsInsideHTMLComment(content, off) {
			idx := next
			if idx > len(content) {
				idx = len(content)
			}
			if strings.HasPrefix(content[idx:], token.NewlineLF) {
				idx += len(token.NewlineLF)
			}
			return idx, true
		}
		off = next
	}
	return 0, false
}

// createConventionSection adds a new section holding the entry, placed
// above "## Themes" so it lands in the staging zone. A root with no
// themes region takes the section at EOF.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted convention bullet
//   - header: the heading line to create
//
// Returns:
//   - []byte: Modified content with the section created
func createConventionSection(content, entry, header string) []byte {
	block := header + token.NewlineLF + token.NewlineLF + entry

	off := 0
	for _, line := range strings.Split(content, token.NewlineLF) {
		if strings.TrimSpace(line) == cfgDisc.HeadingThemes &&
			!ExistsInsideHTMLComment(content, off) {
			return []byte(
				content[:off] + block + token.NewlineLF + content[off:],
			)
		}
		off += len(line) + len(token.NewlineLF)
	}

	if !inspect.EndsWithNewline(content) {
		content += token.NewlineLF
	}
	return []byte(content + token.NewlineLF + block)
}
