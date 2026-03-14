//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/zensical"
)

// GenerateSiteReadme creates a README for the journal-site directory.
//
// Parameters:
//   - journalDir: Path to the source journal directory
//
// Returns:
//   - string: Markdown README content with regeneration instructions
func GenerateSiteReadme(journalDir string) string {
	return fmt.Sprintf(assets.TplJournalSiteReadme, journalDir)
}

// GenerateIndex creates the index.md content for the journal site.
//
// Parameters:
//   - entries: All journal entries to include
//
// Returns:
//   - string: Markdown content for index.md
func GenerateIndex(entries []JournalEntry) string {
	var sb strings.Builder
	nl := token.NewlineLF

	// Separate regular sessions from suggestions and multi-part continuations
	var regular, suggestions []JournalEntry
	for _, e := range entries {
		switch {
		case e.Suggestive:
			suggestions = append(suggestions, e)
		case ContinuesMultipart(e.Filename):
			// Skip part 2+ of split sessions - they're navigable from part 1
			continue
		default:
			regular = append(regular, e)
		}
	}

	sb.WriteString(assets.JournalHeadingSessionJournal + nl + nl)
	sb.WriteString(assets.TplJournalIndexIntro + nl + nl)
	sb.WriteString(fmt.Sprintf(assets.TplJournalIndexStats+
		nl+nl, len(regular), len(suggestions)))

	// Group regular sessions by month
	months, monthOrder := GroupByMonth(regular)

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf(assets.TplJournalMonthHeading+nl+nl, month))

		for _, e := range months[month] {
			sb.WriteString(FormatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	// Suggestions section
	if len(suggestions) > 0 {
		sb.WriteString(token.Separator + nl + nl)
		sb.WriteString(assets.JournalHeadingSuggestions + nl + nl)
		sb.WriteString(assets.TplJournalSuggestionsNote + nl + nl)

		for _, e := range suggestions {
			sb.WriteString(FormatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// FormatIndexEntry formats a single entry for the index.
//
// Parameters:
//   - e: Journal entry to format
//   - nl: Newline string
//
// Returns:
//   - string: Formatted line (e.g., "- 14:30 [title](link.md) (project) `1.2KB`")
func FormatIndexEntry(e JournalEntry, nl string) string {
	link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)

	timeStr := ""
	if e.Time != "" && len(e.Time) >= journal.TimePrefixLen {
		timeStr = e.Time[:journal.TimePrefixLen] + " "
	}

	project := ""
	if e.Project != "" {
		project = fmt.Sprintf(" (%s)", e.Project)
	}

	size := FormatSize(e.Size)

	line := fmt.Sprintf(
		assets.TplJournalIndexEntry+nl, timeStr, e.Title, link, project, size,
	)
	if e.Summary != "" {
		line += fmt.Sprintf(assets.TplJournalIndexSummary+nl, e.Summary)
	}
	return line
}

// InjectSummary inserts the session summary as an admonition after the
// frontmatter (and any source link). Placed before the first heading.
//
// Parameters:
//   - content: Markdown content (may already have source link injected)
//   - summary: Summary text from frontmatter
//
// Returns:
//   - string: Content with the summary admonition injected
func InjectSummary(content, summary string) string {
	nl := token.NewlineLF
	admonition := fmt.Sprintf(
		assets.TplJournalSummaryAdmonition+nl+nl, summary,
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

	// No frontmatter — prepend
	return admonition + content
}

// InjectSourceLink inserts a "View source" link into a journal entry's
// content. The link is placed after YAML frontmatter if present, otherwise
// at the top.
//
// Parameters:
//   - content: Raw Markdown content of the journal entry
//   - sourcePath: Path to the source file on disk
//
// Returns:
//   - string: Content with the source link injected
func InjectSourceLink(content, sourcePath string) string {
	nl := token.NewlineLF
	absPath, pathErr := filepath.Abs(sourcePath)
	if pathErr != nil {
		absPath = sourcePath
	}
	relPath := filepath.Join(
		dir.Context, dir.Journal, filepath.Base(absPath),
	)
	link := fmt.Sprintf(assets.TplJournalSourceLink+nl+nl,
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

// GenerateZensicalToml creates the zensical.toml configuration for the
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
func GenerateZensicalToml(
	entries []JournalEntry, topics []TopicData,
	keyFiles []KeyFileData, sessionTypes []TypeData,
) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(assets.TplZensicalProject + nl)

	// Build navigation
	sb.WriteString(zensical.TomlNavOpen + nl)
	sb.WriteString(fmt.Sprintf(assets.TplJournalNavItem+nl,
		assets.JournalLabelHome, file.Index))
	if len(topics) > 0 {
		sb.WriteString(fmt.Sprintf(assets.TplJournalNavItem+nl,
			assets.JournalLabelTopics,
			filepath.Join(dir.JournTopics, file.Index)),
		)
	}
	if len(keyFiles) > 0 {
		sb.WriteString(fmt.Sprintf(assets.TplJournalNavItem+nl,
			assets.JournalLabelFiles,
			filepath.Join(dir.JournalFiles, file.Index)),
		)
	}
	if len(sessionTypes) > 0 {
		sb.WriteString(fmt.Sprintf(assets.TplJournalNavItem+nl,
			assets.JournalLabelTypes,
			filepath.Join(dir.JournalTypes, file.Index)),
		)
	}

	// Filter out suggestion sessions and multi-part continuations from navigation
	var regular []JournalEntry
	for _, e := range entries {
		if e.Suggestive {
			continue
		}
		if ContinuesMultipart(e.Filename) {
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
		assets.TplJournalNavSection+nl, assets.JournalHeadingRecentSessions),
	)
	for _, e := range recent {
		title := e.Title
		if utf8.RuneCountInString(title) > journal.MaxNavTitleLen {
			runes := []rune(title)
			title = string(runes[:journal.MaxNavTitleLen]) + token.Ellipsis
		}
		title = strings.ReplaceAll(title, `"`, `\"`)
		sb.WriteString(fmt.Sprintf(
			assets.TplJournalNavSessionItem+nl, title, e.Filename),
		)
	}
	sb.WriteString(zensical.TomlNavSectionClose + nl)
	sb.WriteString(zensical.TomlNavClose + nl + nl)

	sb.WriteString(assets.TplZensicalExtraCSS + nl)

	sb.WriteString(assets.TplZensicalTheme)

	return sb.String()
}
