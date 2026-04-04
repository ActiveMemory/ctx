//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package wikilink

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgHTTP "github.com/ActiveMemory/ctx/internal/config/http"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// ConvertMarkdownLinks replaces internal Markdown links with Obsidian
// wikilinks. External links (http/https) are left unchanged.
//
// Parameters:
//   - content: Markdown content with standard links
//
// Returns:
//   - string: Content with internal links converted to [[target|display]]
func ConvertMarkdownLinks(content string) string {
	return regex.MarkdownLinkAny.ReplaceAllStringFunc(
		content, func(match string) string {
			parts := regex.MarkdownLinkAny.FindStringSubmatch(match)
			if len(parts) != 3 {
				return match
			}

			display := parts[1]
			target := parts[2]

			// Skip external links
			if strings.HasPrefix(target, cfgHTTP.PrefixHTTP) ||
				strings.HasPrefix(target, cfgHTTP.PrefixHTTPS) ||
				strings.HasPrefix(target, cfgHTTP.PrefixFile) {
				return match
			}

			// Strip path prefix (e.g., "../topics/", "../") and .md extension
			target = filepath.Base(target)
			target = strings.TrimSuffix(target, file.ExtMarkdown)

			return Format(target, display)
		})
}

// Format formats a wikilink with optional display text.
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
func Format(target, display string) string {
	if target == display {
		return fmt.Sprintf(obsidian.WikilinkPlain, target)
	}
	return fmt.Sprintf(obsidian.WikilinkFmt, target, display)
}

// FormatEntry formats a journal entry as a wikilink list item.
//
// Output: - [[filename|title]] - `type` · `outcome`
//
// Parameters:
//   - e: JournalEntry to format
//
// Returns:
//   - string: Formatted list item with wikilink
func FormatEntry(e entity.JournalEntry) string {
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
		Format(link, e.Title),
		suffix,
	)
}
