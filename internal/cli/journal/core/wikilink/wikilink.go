//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package wikilink

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// RegexMarkdownLink matches Markdown links: [display](target)
var RegexMarkdownLink = regexp.MustCompile(`\[([^]]+)]\(([^)]+)\)`)

// ConvertMarkdownLinks replaces internal Markdown links with Obsidian
// wikilinks. External links (http/https) are left unchanged.
//
// Parameters:
//   - content: Markdown content with standard links
//
// Returns:
//   - string: Content with internal links converted to [[target|display]]
func ConvertMarkdownLinks(content string) string {
	return RegexMarkdownLink.ReplaceAllStringFunc(
		content, func(match string) string {
			parts := RegexMarkdownLink.FindStringSubmatch(match)
			if len(parts) != 3 {
				return match
			}

			display := parts[1]
			target := parts[2]

			// Skip external links
			if strings.HasPrefix(target, "http://") ||
				strings.HasPrefix(target, "https://") ||
				strings.HasPrefix(target, "file://") {
				return match
			}

			// Strip path prefix (e.g., "../topics/", "../") and .md extension
			target = filepath.Base(target)
			target = strings.TrimSuffix(target, file.ExtMarkdown)

			return FormatWikilink(target, display)
		})
}

// FormatWikilink formats a wikilink with optional display text.
//
// If display equals target, a plain wikilink is returned: [[target]]
// Otherwise: [[target|display]]
//
// Parameters:
//   - target: Link target (note name without .md)
//   - display: Display text shown in the link
//
// Returns:
//   - string: Formatted wikilink
func FormatWikilink(target, display string) string {
	if target == display {
		return fmt.Sprintf(obsidian.WikilinkPlain, target)
	}
	return fmt.Sprintf(obsidian.WikilinkFmt, target, display)
}

// FormatWikilinkEntry formats a journal entry as a wikilink list item.
//
// Output: - [[filename|title]] - `type` · `outcome`
//
// Parameters:
//   - e: JournalEntry to format
//
// Returns:
//   - string: Formatted list item with wikilink
func FormatWikilinkEntry(e entity.JournalEntry) string {
	link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)

	var meta []string
	if e.Type != "" {
		meta = append(meta, token.Backtick+e.Type+token.Backtick)
	}
	if e.Outcome != "" {
		meta = append(meta, token.Backtick+e.Outcome+token.Backtick)
	}

	suffix := ""
	if len(meta) > 0 {
		suffix = token.MetaSeparator + strings.Join(meta, token.MetaJoin)
	}

	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteWikilinkListItem),
		FormatWikilink(link, e.Title),
		suffix,
	)
}
