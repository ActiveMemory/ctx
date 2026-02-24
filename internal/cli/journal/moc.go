//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// generateHomeMOC creates the root navigation hub for the Obsidian vault.
//
// The Home MOC links to all section MOCs and lists recent sessions.
//
// Parameters:
//   - entries: All journal entries (filtered, no suggestions/multipart)
//   - hasTopics: Whether any topic data exists
//   - hasFiles: Whether any key file data exists
//   - hasTypes: Whether any type data exists
//
// Returns:
//   - string: Markdown content for Home.md
func generateHomeMOC(
	entries []journalEntry,
	hasTopics, hasFiles, hasTypes bool,
) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString("# Session Journal" + nl + nl)
	sb.WriteString("Navigation hub for all journal entries." + nl + nl)

	sb.WriteString("## Browse by" + nl + nl)
	if hasTopics {
		sb.WriteString(fmt.Sprintf(
			"- %s — sessions grouped by topic"+nl,
			formatWikilink("_Topics", "Topics")))
	}
	if hasFiles {
		sb.WriteString(fmt.Sprintf(
			"- %s — sessions grouped by file touched"+nl,
			formatWikilink("_Key Files", "Key Files")))
	}
	if hasTypes {
		sb.WriteString(fmt.Sprintf(
			"- %s — sessions grouped by type"+nl,
			formatWikilink("_Session Types", "Session Types")))
	}
	sb.WriteString(nl)

	// Recent sessions (up to JournalMaxRecentSessions)
	recent := entries
	if len(recent) > config.JournalMaxRecentSessions {
		recent = recent[:config.JournalMaxRecentSessions]
	}

	sb.WriteString("## Recent Sessions" + nl + nl)
	for _, e := range recent {
		sb.WriteString(formatWikilinkEntry(e) + nl)
	}
	sb.WriteString(nl)

	return sb.String()
}

// generateObsidianTopicsMOC creates the topics index page with wikilinks.
//
// Popular topics link to dedicated pages; long-tail topics link inline
// to the first matching session.
//
// Parameters:
//   - topics: Sorted topic data from buildTopicIndex
//
// Returns:
//   - string: Markdown content for _Topics.md
func generateObsidianTopicsMOC(topics []topicData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	var popular, longtail []topicData
	for _, t := range topics {
		if t.Popular {
			popular = append(popular, t)
		} else {
			longtail = append(longtail, t)
		}
	}

	sb.WriteString("# Topics" + nl + nl)
	sb.WriteString(fmt.Sprintf(
		"**%d topics** across **%d sessions** — **%d popular**, **%d long-tail**"+nl+nl,
		len(topics), countUniqueSessions(topics),
		len(popular), len(longtail)))

	if len(popular) > 0 {
		sb.WriteString("## Popular Topics" + nl + nl)
		for _, t := range popular {
			sb.WriteString(fmt.Sprintf("- %s (%d sessions)"+nl,
				formatWikilink(t.Name, t.Name), len(t.Entries)))
		}
		sb.WriteString(nl)
	}

	if len(longtail) > 0 {
		sb.WriteString("## Long-tail Topics" + nl + nl)
		for _, t := range longtail {
			e := t.Entries[0]
			link := strings.TrimSuffix(e.Filename, config.ExtMarkdown)
			sb.WriteString(fmt.Sprintf("- **%s** — %s"+nl,
				t.Name, formatWikilink(link, e.Title)))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// generateObsidianTopicPage creates an individual topic page with wikilinks
// grouped by month.
//
// Parameters:
//   - topic: Topic data including name and entries
//
// Returns:
//   - string: Markdown content for the topic page
func generateObsidianTopicPage(topic topicData) string {
	return generateObsidianGroupedPage(
		fmt.Sprintf("# %s", topic.Name),
		fmt.Sprintf("**%d sessions** with this topic.", len(topic.Entries)),
		topic.Entries,
	)
}

// generateObsidianFilesMOC creates the key files index page with wikilinks.
//
// Parameters:
//   - keyFiles: Sorted key file data from buildKeyFileIndex
//
// Returns:
//   - string: Markdown content for _Key Files.md
func generateObsidianFilesMOC(keyFiles []keyFileData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	var popular, longtail []keyFileData
	for _, kf := range keyFiles {
		if kf.Popular {
			popular = append(popular, kf)
		} else {
			longtail = append(longtail, kf)
		}
	}

	totalSessions := 0
	seen := make(map[string]bool)
	for _, kf := range keyFiles {
		for _, e := range kf.Entries {
			if !seen[e.Filename] {
				seen[e.Filename] = true
				totalSessions++
			}
		}
	}

	sb.WriteString("# Key Files" + nl + nl)
	sb.WriteString(fmt.Sprintf(
		"**%d files** across **%d sessions** — **%d popular**, **%d long-tail**"+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail)))

	if len(popular) > 0 {
		sb.WriteString("## Frequently Touched" + nl + nl)
		for _, kf := range popular {
			slug := keyFileSlug(kf.Path)
			sb.WriteString(fmt.Sprintf("- %s (%d sessions)"+nl,
				formatWikilink(slug, "`"+kf.Path+"`"),
				len(kf.Entries)))
		}
		sb.WriteString(nl)
	}

	if len(longtail) > 0 {
		sb.WriteString("## Single Session" + nl + nl)
		for _, kf := range longtail {
			e := kf.Entries[0]
			link := strings.TrimSuffix(e.Filename, config.ExtMarkdown)
			sb.WriteString(fmt.Sprintf("- `%s` — %s"+nl,
				kf.Path, formatWikilink(link, e.Title)))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// generateObsidianFilePage creates an individual key file page with wikilinks
// grouped by month.
//
// Parameters:
//   - kf: Key file data including path and entries
//
// Returns:
//   - string: Markdown content for the key file page
func generateObsidianFilePage(kf keyFileData) string {
	return generateObsidianGroupedPage(
		fmt.Sprintf("# `%s`", kf.Path),
		fmt.Sprintf("**%d sessions** touching this file.", len(kf.Entries)),
		kf.Entries,
	)
}

// generateObsidianTypesMOC creates the session types index page with
// wikilinks.
//
// Parameters:
//   - sessionTypes: Sorted type data from buildTypeIndex
//
// Returns:
//   - string: Markdown content for _Session Types.md
func generateObsidianTypesMOC(sessionTypes []typeData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	totalSessions := 0
	for _, st := range sessionTypes {
		totalSessions += len(st.Entries)
	}

	sb.WriteString("# Session Types" + nl + nl)
	sb.WriteString(fmt.Sprintf(
		"**%d types** across **%d sessions**"+nl+nl,
		len(sessionTypes), totalSessions))

	for _, st := range sessionTypes {
		sb.WriteString(fmt.Sprintf("- %s (%d sessions)"+nl,
			formatWikilink(st.Name, st.Name), len(st.Entries)))
	}
	sb.WriteString(nl)

	return sb.String()
}

// generateObsidianTypePage creates an individual session type page with
// wikilinks grouped by month.
//
// Parameters:
//   - st: Type data including name and entries
//
// Returns:
//   - string: Markdown content for the session type page
func generateObsidianTypePage(st typeData) string {
	return generateObsidianGroupedPage(
		fmt.Sprintf("# %s", st.Name),
		fmt.Sprintf("**%d sessions** of type *%s*.", len(st.Entries), st.Name),
		st.Entries,
	)
}

// generateObsidianGroupedPage builds a detail page with a heading, stats line,
// and month-grouped session wikilinks.
//
// Parameters:
//   - heading: Pre-formatted Markdown heading
//   - stats: Pre-formatted stats line
//   - entries: Journal entries to group by month
//
// Returns:
//   - string: Complete Markdown page content
func generateObsidianGroupedPage(
	heading, stats string, entries []journalEntry,
) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(heading + nl + nl)
	sb.WriteString(stats + nl + nl)

	months, monthOrder := groupByMonth(entries)
	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf("## %s"+nl+nl, month))
		for _, e := range months[month] {
			sb.WriteString(formatWikilinkEntry(e) + nl)
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// generateRelatedFooter builds the "Related Sessions" footer appended to
// each journal entry in the vault. Links to topic/type MOCs and lists
// other entries that share topics.
//
// Parameters:
//   - entry: The current journal entry
//   - topicIndex: Map of topic name → entries sharing that topic
//   - maxRelated: Maximum number of "see also" entries to show
//
// Returns:
//   - string: Markdown footer section (empty if entry has no metadata)
func generateRelatedFooter(
	entry journalEntry,
	topicIndex map[string][]journalEntry,
	maxRelated int,
) string {
	if len(entry.Topics) == 0 && entry.Type == "" {
		return ""
	}

	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(nl + config.Separator + nl + nl)
	sb.WriteString(config.ObsidianRelatedHeading + nl + nl)

	// Topic links
	if len(entry.Topics) > 0 {
		topicLinks := make([]string, 0, len(entry.Topics)+1)
		topicLinks = append(topicLinks,
			formatWikilink("_Topics", "Topics MOC"))
		for _, t := range entry.Topics {
			topicLinks = append(topicLinks,
				fmt.Sprintf(config.ObsidianWikilinkPlain, t))
		}
		sb.WriteString("**Topics**: " + strings.Join(topicLinks, " · ") + nl + nl)
	}

	// Type link
	if entry.Type != "" {
		sb.WriteString(fmt.Sprintf("**Type**: %s"+nl+nl,
			fmt.Sprintf(config.ObsidianWikilinkPlain, entry.Type)))
	}

	// See also: other entries sharing topics
	related := collectRelated(entry, topicIndex, maxRelated)
	if len(related) > 0 {
		sb.WriteString(config.ObsidianSeeAlso + nl)
		for _, rel := range related {
			link := strings.TrimSuffix(rel.Filename, config.ExtMarkdown)
			sb.WriteString(fmt.Sprintf("- %s"+nl,
				formatWikilink(link, rel.Title)))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// collectRelated finds entries that share topics with the given entry,
// excluding the entry itself. Returns up to maxRelated unique entries,
// prioritized by number of shared topics.
//
// Parameters:
//   - entry: The current journal entry
//   - topicIndex: Map of topic name → entries
//   - maxRelated: Maximum results
//
// Returns:
//   - []journalEntry: Related entries, deduplicated
func collectRelated(
	entry journalEntry,
	topicIndex map[string][]journalEntry,
	maxRelated int,
) []journalEntry {
	// Count shared topics per entry
	scores := make(map[string]int)
	candidates := make(map[string]journalEntry)

	for _, topic := range entry.Topics {
		for _, rel := range topicIndex[topic] {
			if rel.Filename == entry.Filename {
				continue
			}
			scores[rel.Filename]++
			candidates[rel.Filename] = rel
		}
	}

	// Sort by score descending, then by filename for stability
	type scored struct {
		entry journalEntry
		score int
	}
	var sorted []scored
	for fn, e := range candidates {
		sorted = append(sorted, scored{entry: e, score: scores[fn]})
	}

	// Simple insertion sort (small N)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0; j-- {
			if sorted[j].score > sorted[j-1].score ||
				(sorted[j].score == sorted[j-1].score &&
					sorted[j].entry.Filename < sorted[j-1].entry.Filename) {
				sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
			}
		}
	}

	if len(sorted) > maxRelated {
		sorted = sorted[:maxRelated]
	}

	result := make([]journalEntry, len(sorted))
	for i, s := range sorted {
		result[i] = s.entry
	}
	return result
}
