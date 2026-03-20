//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// GenerateHomeMOC creates the root navigation hub for the Obsidian vault.
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
func GenerateHomeMOC(
	entries []JournalEntry,
	hasTopics, hasFiles, hasTypes bool,
) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(desc.TextDesc(text.DescKeyHeadingSessionJournal) + nl + nl)
	sb.WriteString(desc.TextDesc(text.DescKeyJournalMocNavDescription) + nl + nl)

	sb.WriteString(desc.TextDesc(text.DescKeyJournalMocBrowseBy) + nl + nl)
	if hasTopics {
		sb.WriteString(fmt.Sprintf(
			"- %s %s"+nl,
			FormatWikilink("_Topics", "Topics"),
			desc.TextDesc(text.DescKeyJournalMocTopicsDesc)))
	}
	if hasFiles {
		sb.WriteString(fmt.Sprintf(
			"- %s %s"+nl,
			FormatWikilink("_Key Files", "Key Files"),
			desc.TextDesc(text.DescKeyJournalMocFilesDesc)))
	}
	if hasTypes {
		sb.WriteString(fmt.Sprintf(
			"- %s %s"+nl,
			FormatWikilink("_Session Types", "Session Types"),
			desc.TextDesc(text.DescKeyJournalMocTypesDesc)))
	}
	sb.WriteString(nl)

	// Recent sessions (up to MaxRecentSessions)
	recent := entries
	if len(recent) > journal.MaxRecentSessions {
		recent = recent[:journal.MaxRecentSessions]
	}

	sb.WriteString(token.HeadingLevelTwoStart + desc.TextDesc(text.DescKeyHeadingRecentSessions) + nl + nl)
	for _, e := range recent {
		sb.WriteString(FormatWikilinkEntry(e) + nl)
	}
	sb.WriteString(nl)

	return sb.String()
}

// GenerateObsidianTopicsMOC creates the topics index page with wikilinks.
//
// Popular topics link to dedicated pages; long-tail topics link inline
// to the first matching session.
//
// Parameters:
//   - topics: Sorted topic data from BuildTopicIndex
//
// Returns:
//   - string: Markdown content for _Topics.md
func GenerateObsidianTopicsMOC(topics []TopicData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	var popular, longtail []TopicData
	for _, t := range topics {
		if t.Popular {
			popular = append(popular, t)
		} else {
			longtail = append(longtail, t)
		}
	}

	sb.WriteString("# Topics" + nl + nl)
	sb.WriteString(fmt.Sprintf(
		desc.TextDesc(text.DescKeyJournalMocTopicStats)+nl+nl,
		len(topics), CountUniqueSessions(topics),
		len(popular), len(longtail)))

	if len(popular) > 0 {
		sb.WriteString("## Popular Topics" + nl + nl)
		for _, t := range popular {
			sb.WriteString(fmt.Sprintf("- %s (%d sessions)"+nl,
				FormatWikilink(t.Name, t.Name), len(t.Entries)))
		}
		sb.WriteString(nl)
	}

	if len(longtail) > 0 {
		sb.WriteString("## Long-tail Topics" + nl + nl)
		for _, t := range longtail {
			e := t.Entries[0]
			link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)
			sb.WriteString(fmt.Sprintf("- **%s** — %s"+nl,
				t.Name, FormatWikilink(link, e.Title)))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// GenerateObsidianTopicPage creates an individual topic page with wikilinks
// grouped by month.
//
// Parameters:
//   - topic: Topic data including name and entries
//
// Returns:
//   - string: Markdown content for the topic page
func GenerateObsidianTopicPage(topic TopicData) string {
	return GenerateObsidianGroupedPage(
		fmt.Sprintf("# %s", topic.Name),
		fmt.Sprintf(desc.TextDesc(text.DescKeyJournalMocTopicPageStats), len(topic.Entries)),
		topic.Entries,
	)
}

// GenerateObsidianFilesMOC creates the key files index page with wikilinks.
//
// Parameters:
//   - keyFiles: Sorted key file data from BuildKeyFileIndex
//
// Returns:
//   - string: Markdown content for _Key Files.md
func GenerateObsidianFilesMOC(keyFiles []KeyFileData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	var popular, longtail []KeyFileData
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
		desc.TextDesc(text.DescKeyJournalMocFileStats)+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail)))

	if len(popular) > 0 {
		sb.WriteString("## Frequently Touched" + nl + nl)
		for _, kf := range popular {
			slug := KeyFileSlug(kf.Path)
			sb.WriteString(fmt.Sprintf("- %s (%d sessions)"+nl,
				FormatWikilink(slug, "`"+kf.Path+"`"),
				len(kf.Entries)))
		}
		sb.WriteString(nl)
	}

	if len(longtail) > 0 {
		sb.WriteString("## Single Session" + nl + nl)
		for _, kf := range longtail {
			e := kf.Entries[0]
			link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)
			sb.WriteString(fmt.Sprintf("- `%s` — %s"+nl,
				kf.Path, FormatWikilink(link, e.Title)))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// GenerateObsidianFilePage creates an individual key file page with wikilinks
// grouped by month.
//
// Parameters:
//   - kf: Key file data including path and entries
//
// Returns:
//   - string: Markdown content for the key file page
func GenerateObsidianFilePage(kf KeyFileData) string {
	return GenerateObsidianGroupedPage(
		fmt.Sprintf("# `%s`", kf.Path),
		fmt.Sprintf(desc.TextDesc(text.DescKeyJournalMocFilePageStats), len(kf.Entries)),
		kf.Entries,
	)
}

// GenerateObsidianTypesMOC creates the session types index page with
// wikilinks.
//
// Parameters:
//   - sessionTypes: Sorted type data from BuildTypeIndex
//
// Returns:
//   - string: Markdown content for _Session Types.md
func GenerateObsidianTypesMOC(sessionTypes []TypeData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	totalSessions := 0
	for _, st := range sessionTypes {
		totalSessions += len(st.Entries)
	}

	sb.WriteString("# Session Types" + nl + nl)
	sb.WriteString(fmt.Sprintf(
		desc.TextDesc(text.DescKeyJournalMocTypeStats)+nl+nl,
		len(sessionTypes), totalSessions))

	for _, st := range sessionTypes {
		sb.WriteString(fmt.Sprintf("- %s (%d sessions)"+nl,
			FormatWikilink(st.Name, st.Name), len(st.Entries)))
	}
	sb.WriteString(nl)

	return sb.String()
}

// GenerateObsidianTypePage creates an individual session type page with
// wikilinks grouped by month.
//
// Parameters:
//   - st: Type data including name and entries
//
// Returns:
//   - string: Markdown content for the session type page
func GenerateObsidianTypePage(st TypeData) string {
	return GenerateObsidianGroupedPage(
		fmt.Sprintf("# %s", st.Name),
		fmt.Sprintf(desc.TextDesc(text.DescKeyJournalMocTypePageStats), len(st.Entries), st.Name),
		st.Entries,
	)
}

// GenerateObsidianGroupedPage builds a detail page with a heading, stats line,
// and month-grouped session wikilinks.
//
// Parameters:
//   - heading: Pre-formatted Markdown heading
//   - stats: Pre-formatted stats line
//   - entries: Journal entries to group by month
//
// Returns:
//   - string: Complete Markdown page content
func GenerateObsidianGroupedPage(
	heading, stats string, entries []JournalEntry,
) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(heading + nl + nl)
	sb.WriteString(stats + nl + nl)

	months, monthOrder := GroupByMonth(entries)
	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf("## %s"+nl+nl, month))
		for _, e := range months[month] {
			sb.WriteString(FormatWikilinkEntry(e) + nl)
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// GenerateRelatedFooter builds the "Related Sessions" footer appended to
// each journal entry in the vault. Links to topic/type MOCs and lists
// other entries that share topics.
//
// Parameters:
//   - entry: The current journal entry
//   - topicIndex: Map of topic name -> entries sharing that topic
//   - maxRelated: Maximum number of "see also" entries to show
//
// Returns:
//   - string: Markdown footer section (empty if entry has no metadata)
func GenerateRelatedFooter(
	entry JournalEntry,
	topicIndex map[string][]JournalEntry,
	maxRelated int,
) string {
	if len(entry.Topics) == 0 && entry.Type == "" {
		return ""
	}

	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(nl + token.Separator + nl + nl)
	sb.WriteString(desc.TextDesc(text.DescKeyHeadingObsidianRelated) + nl + nl)

	// Topic links
	if len(entry.Topics) > 0 {
		topicLinks := make([]string, 0, len(entry.Topics)+1)
		topicLinks = append(topicLinks,
			FormatWikilink("_Topics", "Topics MOC"))
		for _, t := range entry.Topics {
			topicLinks = append(topicLinks,
				fmt.Sprintf(obsidian.WikilinkPlain, t))
		}
		sb.WriteString(desc.TextDesc(text.DescKeyJournalMocTopicsLabel) + strings.Join(topicLinks, " · ") + nl + nl)
	}

	// Type link
	if entry.Type != "" {
		sb.WriteString(fmt.Sprintf(desc.TextDesc(text.DescKeyJournalMocTypeLabel)+"%s"+nl+nl,
			fmt.Sprintf(obsidian.WikilinkPlain, entry.Type)))
	}

	// See also: other entries sharing topics
	related := CollectRelated(entry, topicIndex, maxRelated)
	if len(related) > 0 {
		sb.WriteString(desc.TextDesc(text.DescKeyLabelObsidianSeeAlso) + nl)
		for _, rel := range related {
			link := strings.TrimSuffix(rel.Filename, file.ExtMarkdown)
			sb.WriteString(fmt.Sprintf("- %s"+nl,
				FormatWikilink(link, rel.Title)))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// CollectRelated finds entries that share topics with the given entry,
// excluding the entry itself. Returns up to maxRelated unique entries,
// prioritized by number of shared topics.
//
// Parameters:
//   - entry: The current journal entry
//   - topicIndex: Map of topic name -> entries
//   - maxRelated: Maximum results
//
// Returns:
//   - []JournalEntry: Related entries, deduplicated
func CollectRelated(
	entry JournalEntry,
	topicIndex map[string][]JournalEntry,
	maxRelated int,
) []JournalEntry {
	// Count shared topics per entry
	scores := make(map[string]int)
	candidates := make(map[string]JournalEntry)

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
		entry JournalEntry
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

	result := make([]JournalEntry, len(sorted))
	for i, s := range sorted {
		result[i] = s.entry
	}
	return result
}

// FilterRegularEntries returns entries excluding suggestions and multipart
// continuations.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []JournalEntry: Filtered entries
func FilterRegularEntries(entries []JournalEntry) []JournalEntry {
	var result []JournalEntry
	for _, e := range entries {
		if e.Suggestive || ContinuesMultipart(e.Filename) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// FilterEntriesWithTopics returns non-suggestive, non-multipart entries
// that have topics.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []JournalEntry: Entries with topics
func FilterEntriesWithTopics(entries []JournalEntry) []JournalEntry {
	var result []JournalEntry
	for _, e := range entries {
		if e.Suggestive || ContinuesMultipart(e.Filename) || len(e.Topics) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

// FilterEntriesWithKeyFiles returns non-suggestive, non-multipart entries
// that have key files.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []JournalEntry: Entries with key files
func FilterEntriesWithKeyFiles(entries []JournalEntry) []JournalEntry {
	var result []JournalEntry
	for _, e := range entries {
		if e.Suggestive || ContinuesMultipart(e.Filename) || len(e.KeyFiles) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

// FilterEntriesWithType returns non-suggestive, non-multipart entries
// that have a type.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []JournalEntry: Entries with type
func FilterEntriesWithType(entries []JournalEntry) []JournalEntry {
	var result []JournalEntry
	for _, e := range entries {
		if e.Suggestive || ContinuesMultipart(e.Filename) || e.Type == "" {
			continue
		}
		result = append(result, e)
	}
	return result
}

// BuildTopicLookup creates a map from topic name to all entries with
// that topic, for efficient related-entry lookups.
//
// Parameters:
//   - entries: Entries with topics
//
// Returns:
//   - map[string][]JournalEntry: Topic name -> entries
func BuildTopicLookup(entries []JournalEntry) map[string][]JournalEntry {
	lookup := make(map[string][]JournalEntry)
	for _, e := range entries {
		for _, topic := range e.Topics {
			lookup[topic] = append(lookup[topic], e)
		}
	}
	return lookup
}
