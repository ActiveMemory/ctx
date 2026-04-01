//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package generate

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/group"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/section"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/zensical"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// SiteReadme creates a README for the journal-site directory.
//
// Parameters:
//   - journalDir: Path to the source journal directory
//
// Returns:
//   - string: Markdown README content with regeneration instructions
func SiteReadme(journalDir string) string {
	return fmt.Sprintf(tpl.JournalSiteReadme, journalDir)
}

// Index creates the index.md content for the journal site.
//
// Parameters:
//   - entries: All journal entries to include
//
// Returns:
//   - string: Markdown content for index.md
func Index(entries []entity.JournalEntry) string {
	var sb strings.Builder
	nl := token.NewlineLF

	// Separate regular sessions from suggestions and multi-part continuations
	var regular, suggestions []entity.JournalEntry
	for _, e := range entries {
		switch {
		case e.Suggestive:
			suggestions = append(suggestions, e)
		case section.ContinuesMultipart(e.Filename):
			// Skip part 2+ of split sessions - they're navigable from part 1
			continue
		default:
			regular = append(regular, e)
		}
	}

	sb.WriteString(desc.Text(text.DescKeyHeadingSessionJournal) + nl + nl)
	sb.WriteString(tpl.JournalIndexIntro + nl + nl)
	sb.WriteString(fmt.Sprintf(tpl.JournalIndexStats+
		nl+nl, len(regular), len(suggestions)))

	// Group regular sessions by month
	months, monthOrder := group.ByMonth(regular)

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf(tpl.JournalMonthHeading+nl+nl, month))

		for _, e := range months[month] {
			sb.WriteString(formatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	// Suggestions section
	if len(suggestions) > 0 {
		sb.WriteString(token.Separator + nl + nl)
		sb.WriteString(desc.Text(text.DescKeyHeadingSuggestions) + nl + nl)
		sb.WriteString(tpl.JournalSuggestionsNote + nl + nl)

		for _, e := range suggestions {
			sb.WriteString(formatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// InjectedSummary inserts the session summary as an admonition after the
// frontmatter (and any source link). Placed before the first heading.
//
// Parameters:
//   - content: Markdown content (may already have the source link injected)
//   - summary: InjectedSummary text from frontmatter
//
// Returns:
//   - string: Content with the summary admonition injected
func InjectedSummary(content, summary string) string {
	nl := token.NewlineLF
	admonition := fmt.Sprintf(
		tpl.JournalSummaryAdmonition+nl+nl, summary,
	)

	// Insert after frontmatter closing delimiter
	fmOpen := len(token.Separator + nl)
	fmClose := len(nl + token.Separator + nl)
	if strings.HasPrefix(content, token.Separator+nl) {
		if end := strings.Index(content[fmOpen:], nl+
			token.Separator+nl); end >= 0 {
			insertAt := fmOpen + end + fmClose
			// Skip past any existing blank lines + source link after frontmatter
			rest := content[insertAt:]
			return content[:insertAt] + nl + admonition + rest
		}
	}

	// No frontmatter - prepend
	return admonition + content
}

// InjectedSourceLink inserts a "View source" link into a journal entry's
// content. The link is placed after YAML frontmatter if present, otherwise
// at the top.
//
// Parameters:
//   - content: Raw Markdown content of the journal entry
//   - sourcePath: Path to the source file on disk
//
// Returns:
//   - string: Content with the source link injected
func InjectedSourceLink(content, sourcePath string) string {
	nl := token.NewlineLF
	absPath, pathErr := filepath.Abs(sourcePath)
	if pathErr != nil {
		absPath = sourcePath
	}
	relPath := filepath.Join(
		dir.Context, dir.Journal, filepath.Base(absPath),
	)
	link := fmt.Sprintf(tpl.JournalSourceLink+nl+nl,
		absPath, relPath, relPath)

	fmOpen := len(token.Separator + nl)
	fmClose := len(nl + token.Separator + nl)
	if strings.HasPrefix(content, token.Separator+nl) {
		if end := strings.Index(content[fmOpen:], nl+
			token.Separator+nl); end >= 0 {
			insertAt := fmOpen + end + fmClose
			return content[:insertAt] + nl + link + content[insertAt:]
		}
	}

	return link + content
}

// ZensicalToml creates the zensical.toml configuration for the
// journal site.
//
// Parameters:
//   - entries: All journal entries for navigation
//   - topics: Topic index data for nav links
//   - keyFiles: Key file index data for nav links
//   - sessionTypes: Session type index data for nav links
//
// Returns:
//   - string: Complete zensical.toml content
func ZensicalToml(
	entries []entity.JournalEntry, topics []entity.TopicData,
	keyFiles []entity.KeyFileData, sessionTypes []entity.TypeData,
) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(tpl.ZensicalProject + nl)

	// Build navigation
	sb.WriteString(zensical.TomlNavOpen + nl)
	sb.WriteString(fmt.Sprintf(tpl.JournalNavItem+nl,
		desc.Text(text.DescKeyLabelHome), file.Index))
	if len(topics) > 0 {
		sb.WriteString(fmt.Sprintf(tpl.JournalNavItem+nl,
			desc.Text(text.DescKeyLabelTopics),
			filepath.Join(dir.JournTopics, file.Index)),
		)
	}
	if len(keyFiles) > 0 {
		sb.WriteString(fmt.Sprintf(tpl.JournalNavItem+nl,
			desc.Text(text.DescKeyLabelFiles),
			filepath.Join(dir.JournalFiles, file.Index)),
		)
	}
	if len(sessionTypes) > 0 {
		sb.WriteString(fmt.Sprintf(tpl.JournalNavItem+nl,
			desc.Text(text.DescKeyLabelTypes),
			filepath.Join(dir.JournalTypes, file.Index)),
		)
	}

	// Filter out suggestion sessions and multi-part continuations from navigation
	var regular []entity.JournalEntry
	for _, e := range entries {
		if e.Suggestive {
			continue
		}
		if section.ContinuesMultipart(e.Filename) {
			continue
		}
		regular = append(regular, e)
	}

	// Group recent entries (last N, excluding suggestions)
	recent := regular
	if len(recent) > journal.MaxRecentSessions {
		recent = recent[:journal.MaxRecentSessions]
	}

	sb.WriteString(fmt.Sprintf(
		tpl.JournalNavSection+nl, desc.Text(text.DescKeyHeadingRecentSessions)),
	)
	for _, e := range recent {
		title := e.Title
		if utf8.RuneCountInString(title) > journal.MaxNavTitleLen {
			runes := []rune(title)
			title = string(runes[:journal.MaxNavTitleLen]) + token.Ellipsis
		}
		title = strings.ReplaceAll(title, token.DoubleQuote, token.EscapedDoubleQuote)
		sb.WriteString(fmt.Sprintf(
			tpl.JournalNavSessionItem+nl, title, e.Filename),
		)
	}
	sb.WriteString(zensical.TomlNavSectionClose + nl)
	sb.WriteString(zensical.TomlNavClose + nl + nl)

	sb.WriteString(tpl.ZensicalExtraCSS + nl)

	sb.WriteString(tpl.ZensicalTheme)

	return sb.String()
}
