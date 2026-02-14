//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/config"
)

// generateSiteReadme creates a README for the journal-site directory.
//
// Parameters:
//   - journalDir: Path to the source journal directory
//
// Returns:
//   - string: Markdown README content with regeneration instructions
func generateSiteReadme(journalDir string) string {
	return fmt.Sprintf(config.TplJournalSiteReadme, journalDir)
}

// generateIndex creates the index.md content for the journal site.
//
// Parameters:
//   - entries: All journal entries to include
//
// Returns:
//   - string: Markdown content for index.md
func generateIndex(entries []journalEntry) string {
	var sb strings.Builder
	nl := config.NewlineLF

	// Separate regular sessions from suggestions and multi-part continuations
	var regular, suggestions []journalEntry
	for _, e := range entries {
		switch {
		case e.Suggestive:
			suggestions = append(suggestions, e)
		case continuesMultipart(e.Filename):
			// Skip part 2+ of split sessions - they're navigable from part 1
			continue
		default:
			regular = append(regular, e)
		}
	}

	sb.WriteString(config.JournalHeadingSessionJournal + nl + nl)
	sb.WriteString(config.TplJournalIndexIntro + nl + nl)
	sb.WriteString(fmt.Sprintf(config.TplJournalIndexStats+
		nl+nl, len(regular), len(suggestions)))

	// Group regular sessions by month
	months, monthOrder := groupByMonth(regular)

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf(config.TplJournalMonthHeading+nl+nl, month))

		for _, e := range months[month] {
			sb.WriteString(formatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	// Suggestions section
	if len(suggestions) > 0 {
		sb.WriteString(config.Separator + nl + nl)
		sb.WriteString(config.JournalHeadingSuggestions + nl + nl)
		sb.WriteString(config.TplJournalSuggestionsNote + nl + nl)

		for _, e := range suggestions {
			sb.WriteString(formatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// formatIndexEntry formats a single entry for the index.
//
// Parameters:
//   - e: Journal entry to format
//   - nl: Newline string
//
// Returns:
//   - string: Formatted line (e.g., "- 14:30 [title](link.md) (project) `1.2KB`")
func formatIndexEntry(e journalEntry, nl string) string {
	link := strings.TrimSuffix(e.Filename, config.ExtMarkdown)

	timeStr := ""
	if e.Time != "" && len(e.Time) >= config.JournalTimePrefixLen {
		timeStr = e.Time[:config.JournalTimePrefixLen] + " "
	}

	project := ""
	if e.Project != "" {
		project = fmt.Sprintf(" (%s)", e.Project)
	}

	size := formatSize(e.Size)

	return fmt.Sprintf(
		config.TplJournalIndexEntry+nl, timeStr, e.Title, link, project, size,
	)
}

// injectSourceLink inserts a "View source" link into a journal entry's
// content. The link is placed after YAML frontmatter if present, otherwise
// at the top.
//
// Parameters:
//   - content: Raw Markdown content of the journal entry
//   - sourcePath: Path to the source file on disk
//
// Returns:
//   - string: Content with the source link injected
func injectSourceLink(content, sourcePath string) string {
	nl := config.NewlineLF
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		absPath = sourcePath
	}
	relPath := filepath.Join(
		config.DirContext, config.DirJournal, filepath.Base(absPath),
	)
	link := fmt.Sprintf(config.TplJournalSourceLink+nl+nl,
		absPath, relPath, relPath)

	fmOpen := len(config.Separator + nl)
	fmClose := len(nl + config.Separator + nl)
	if strings.HasPrefix(content, config.Separator+nl) {
		if end := strings.Index(content[fmOpen:], nl+
			config.Separator+nl); end >= 0 {
			insertAt := fmOpen + end + fmClose
			return content[:insertAt] + nl + link + content[insertAt:]
		}
	}

	return link + content
}

// generateZensicalToml creates the zensical.toml configuration for the
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
func generateZensicalToml(
	entries []journalEntry, topics []topicData,
	keyFiles []keyFileData, sessionTypes []typeData,
) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(config.TplZensicalProject + nl)

	// Build navigation
	sb.WriteString(config.TomlNavOpen + nl)
	sb.WriteString(fmt.Sprintf(config.TplJournalNavItem+nl,
		config.JournalLabelHome, config.FilenameIndex))
	if len(topics) > 0 {
		sb.WriteString(fmt.Sprintf(config.TplJournalNavItem+nl,
			config.JournalLabelTopics,
			filepath.Join(config.JournalDirTopics, config.FilenameIndex)),
		)
	}
	if len(keyFiles) > 0 {
		sb.WriteString(fmt.Sprintf(config.TplJournalNavItem+nl,
			config.JournalLabelFiles,
			filepath.Join(config.JournalDirFiles, config.FilenameIndex)),
		)
	}
	if len(sessionTypes) > 0 {
		sb.WriteString(fmt.Sprintf(config.TplJournalNavItem+nl,
			config.JournalLabelTypes,
			filepath.Join(config.JournalDirTypes, config.FilenameIndex)),
		)
	}

	// Filter out suggestion sessions and multi-part continuations from navigation
	var regular []journalEntry
	for _, e := range entries {
		if e.Suggestive {
			continue
		}
		if continuesMultipart(e.Filename) {
			continue
		}
		regular = append(regular, e)
	}

	// Group recent entries (last N, excluding suggestions)
	recent := regular
	if len(recent) > config.JournalMaxRecentSessions {
		recent = recent[:config.JournalMaxRecentSessions]
	}

	sb.WriteString(fmt.Sprintf(
		config.TplJournalNavSection+nl, config.JournalHeadingRecentSessions),
	)
	for _, e := range recent {
		title := e.Title
		if utf8.RuneCountInString(title) > config.JournalMaxNavTitleLen {
			runes := []rune(title)
			title = string(runes[:config.JournalMaxNavTitleLen]) + config.Ellipsis
		}
		title = strings.ReplaceAll(title, `"`, `\"`)
		sb.WriteString(fmt.Sprintf(
			config.TplJournalNavSessionItem+nl, title, e.Filename),
		)
	}
	sb.WriteString(config.TomlNavSectionClose + nl)
	sb.WriteString(config.TomlNavClose + nl + nl)

	sb.WriteString(config.TplZensicalTheme)

	return sb.String()
}
