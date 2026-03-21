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

	sb.WriteString(desc.Text(text.DescKeyHeadingSessionJournal) + nl + nl)
	sb.WriteString(desc.Text(text.DescKeyJournalMocNavDescription) + nl + nl)

	sb.WriteString(desc.Text(text.DescKeyJournalMocBrowseBy) + nl + nl)
	browseItem := desc.Text(text.DescKeyJournalMocBrowseItem)
	topicsLink := strings.TrimSuffix(obsidian.MOCTopics, file.ExtMarkdown)
	filesLink := strings.TrimSuffix(obsidian.MOCFiles, file.ExtMarkdown)
	typesLink := strings.TrimSuffix(obsidian.MOCTypes, file.ExtMarkdown)
	if hasTopics {
		sb.WriteString(fmt.Sprintf(
			browseItem+nl,
			FormatWikilink(topicsLink, topicsLink[1:]),
			desc.Text(text.DescKeyJournalMocTopicsDesc)))
	}
	if hasFiles {
		sb.WriteString(fmt.Sprintf(
			browseItem+nl,
			FormatWikilink(filesLink, filesLink[1:]),
			desc.Text(text.DescKeyJournalMocFilesDesc)))
	}
	if hasTypes {
		sb.WriteString(fmt.Sprintf(
			browseItem+nl,
			FormatWikilink(typesLink, typesLink[1:]),
			desc.Text(text.DescKeyJournalMocTypesDesc)))
	}
	sb.WriteString(nl)

	// Recent sessions (up to MaxRecentSessions)
	recent := entries
	if len(recent) > journal.MaxRecentSessions {
		recent = recent[:journal.MaxRecentSessions]
	}

	sb.WriteString(
		token.HeadingLevelTwoStart +
			desc.Text(text.DescKeyHeadingRecentSessions) + nl + nl,
	)
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

	popular, longtail := SplitPopular(topics)

	sb.WriteString(desc.Text(text.DescKeyJournalMocHeadingTopics) + nl + nl)
	sb.WriteString(fmt.Sprintf(
		desc.Text(text.DescKeyJournalMocTopicStats)+nl+nl,
		len(topics), CountUniqueSessions(topics),
		len(popular), len(longtail)))

	writeSection(
		&sb,
		text.DescKeyJournalMocHeadingPopular, popular, func(t TopicData) string {
			return fmt.Sprintf(
				desc.Text(text.DescKeyJournalMocItemSessions)+nl,
				FormatWikilink(t.Name, t.Name), len(t.Entries))
		})
	writeSection(
		&sb,
		text.DescKeyJournalMocHeadingLongtail,
		longtail, func(t TopicData) string {
			e := t.Entries[0]
			link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)
			return fmt.Sprintf(
				desc.Text(text.DescKeyJournalMocItemNamed)+nl,
				t.Name, FormatWikilink(link, e.Title))
		})

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
		fmt.Sprintf(desc.Text(text.DescKeyJournalMocPageTitle), topic.Name),
		fmt.Sprintf(
			desc.Text(text.DescKeyJournalMocTopicPageStats),
			len(topic.Entries),
		),
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

	popular, longtail := SplitPopular(keyFiles)

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

	sb.WriteString(desc.Text(text.DescKeyJournalMocHeadingFiles) + nl + nl)
	sb.WriteString(fmt.Sprintf(
		desc.Text(text.DescKeyJournalMocFileStats)+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail)))

	writeSection(
		&sb,
		text.DescKeyJournalMocHeadingFreq, popular,
		func(kf KeyFileData) string {
			slug := KeyFileSlug(kf.Path)
			return fmt.Sprintf(
				desc.Text(text.DescKeyJournalMocItemFileSess)+nl,
				FormatWikilink(slug, "`"+kf.Path+"`"),
				len(kf.Entries))
		})
	writeSection(
		&sb,
		text.DescKeyJournalMocHeadingSingle,
		longtail, func(kf KeyFileData) string {
			e := kf.Entries[0]
			link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)
			return fmt.Sprintf(
				desc.Text(text.DescKeyJournalMocItemFileNamed)+nl,
				kf.Path, FormatWikilink(link, e.Title))
		})

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
		fmt.Sprintf(desc.Text(text.DescKeyJournalMocCodeTitle), kf.Path),
		fmt.Sprintf(
			desc.Text(text.DescKeyJournalMocFilePageStats), len(kf.Entries),
		),
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

	sb.WriteString(desc.Text(text.DescKeyJournalMocHeadingTypes) + nl + nl)
	sb.WriteString(fmt.Sprintf(
		desc.Text(text.DescKeyJournalMocTypeStats)+nl+nl,
		len(sessionTypes), totalSessions))

	for _, st := range sessionTypes {
		sb.WriteString(fmt.Sprintf(
			desc.Text(text.DescKeyJournalMocItemSessions)+nl,
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
		fmt.Sprintf(desc.Text(text.DescKeyJournalMocPageTitle), st.Name),
		fmt.Sprintf(desc.Text(text.DescKeyJournalMocTypePageStats), len(st.Entries), st.Name),
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
		sb.WriteString(fmt.Sprintf(desc.Text(text.DescKeyJournalMocHeadingMonth)+nl+nl, month))
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
	sb.WriteString(desc.Text(text.DescKeyHeadingObsidianRelated) + nl + nl)

	// Topic links
	if len(entry.Topics) > 0 {
		topicLinks := make([]string, 0, len(entry.Topics)+1)
		topicLinks = append(topicLinks,
			FormatWikilink(
				strings.TrimSuffix(obsidian.MOCTopics, file.ExtMarkdown),
				desc.Text(text.DescKeyJournalMocTopicsMocLink)))
		for _, t := range entry.Topics {
			topicLinks = append(topicLinks,
				fmt.Sprintf(obsidian.WikilinkPlain, t))
		}
		sb.WriteString(fmt.Sprintf(
			desc.Text(text.DescKeyJournalMocTopicsLabel)+nl+nl,
			strings.Join(topicLinks, desc.Text(text.DescKeyJournalMocTopicSep))))
	}

	// Type link
	if entry.Type != "" {
		sb.WriteString(fmt.Sprintf(
			desc.Text(text.DescKeyJournalMocTypeLabel)+nl+nl,
			fmt.Sprintf(obsidian.WikilinkPlain, entry.Type)))
	}

	// See also: other entries sharing topics
	related := CollectRelated(entry, topicIndex, maxRelated)
	if len(related) > 0 {
		sb.WriteString(desc.Text(text.DescKeyLabelObsidianSeeAlso) + nl)
		for _, rel := range related {
			link := strings.TrimSuffix(rel.Filename, file.ExtMarkdown)
			sb.WriteString(fmt.Sprintf(
				desc.Text(text.DescKeyJournalMocItemListed)+nl,
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

// BuildTopicLookup creates a map from the topic name to all entries with
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
