//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// regexMarkdownLink matches Markdown links: [display](target)
var regexMarkdownLink = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

// convertMarkdownLinks replaces internal Markdown links with Obsidian
// wikilinks. External links (http/https) are left unchanged.
//
// Parameters:
//   - content: Markdown content with standard links
//
// Returns:
//   - string: Content with internal links converted to [[target|display]]
func convertMarkdownLinks(content string) string {
	return regexMarkdownLink.ReplaceAllStringFunc(content, func(match string) string {
		parts := regexMarkdownLink.FindStringSubmatch(match)
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
		target = strings.TrimSuffix(target, config.ExtMarkdown)

		return formatWikilink(target, display)
	})
}

// formatWikilink formats a wikilink with optional display text.
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
func formatWikilink(target, display string) string {
	if target == display {
		return fmt.Sprintf(config.ObsidianWikilinkPlain, target)
	}
	return fmt.Sprintf(config.ObsidianWikilinkFmt, target, display)
}

// formatWikilinkEntry formats a journal entry as a wikilink list item.
//
// Output: - [[filename|title]] — `type` · `outcome`
//
// Parameters:
//   - e: Journal entry to format
//
// Returns:
//   - string: Formatted list item with wikilink
func formatWikilinkEntry(e journalEntry) string {
	link := strings.TrimSuffix(e.Filename, config.ExtMarkdown)

	var meta []string
	if e.Type != "" {
		meta = append(meta, config.Backtick+e.Type+config.Backtick)
	}
	if e.Outcome != "" {
		meta = append(meta, config.Backtick+e.Outcome+config.Backtick)
	}

	suffix := ""
	if len(meta) > 0 {
		suffix = " — " + strings.Join(meta, " · ")
	}

	return fmt.Sprintf("- %s%s",
		formatWikilink(link, e.Title), suffix)
}
