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
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// BuildTopicIndex aggregates entries by topic and returns sorted topic data.
//
// Topics with 2+ sessions are marked popular. Sorted by count desc, then alpha.
//
// Parameters:
//   - entries: All journal entries to aggregate
//
// Returns:
//   - []TopicData: Topics sorted by session count descending, then name
func BuildTopicIndex(entries []JournalEntry) []TopicData {
	grouped := BuildGroupedIndex(
		entries,
		func(e JournalEntry) []string { return e.Topics },
	)
	topics := make([]TopicData, len(grouped))
	for i, g := range grouped {
		topics[i] = TopicData{Name: g.Key, Entries: g.Entries, Popular: g.Popular}
	}
	return topics
}

// GenerateTopicsIndex creates the topics/index.md page.
//
// Popular topics link to dedicated pages; long-tail topics list entries inline.
//
// Parameters:
//   - topics: Sorted topic data from BuildTopicIndex
//
// Returns:
//   - string: Markdown content for topics/index.md
func GenerateTopicsIndex(topics []TopicData) string {
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

	sb.WriteString(desc.TextDesc(text.DescKeyHeadingTopics) + nl + nl)
	sb.WriteString(fmt.Sprintf(
		tpl.TplJournalTopicStats+nl+nl,
		len(topics), CountUniqueSessions(topics), len(popular), len(longtail)))

	WritePopularAndLongtail(&sb,
		len(popular), desc.TextDesc(text.DescKeyHeadingPopularTopics),
		func(i int) (string, string, int) {
			return popular[i].Name, popular[i].Name, len(popular[i].Entries)
		},
		len(longtail), desc.TextDesc(text.DescKeyHeadingLongtailTopics),
		tpl.TplJournalLongtailEntry,
		func(i int) (string, JournalEntry) {
			return longtail[i].Name, longtail[i].Entries[0]
		},
	)

	return sb.String()
}

// GenerateTopicPage creates an individual topic page with sessions grouped
// by month.
//
// Parameters:
//   - topic: Topic data including name and entries
//
// Returns:
//   - string: Markdown content for the topic page
func GenerateTopicPage(topic TopicData) string {
	return GenerateGroupedPage(
		fmt.Sprintf(tpl.TplJournalPageHeading, topic.Name),
		fmt.Sprintf(tpl.TplJournalTopicPageStats, len(topic.Entries)),
		topic.Entries,
	)
}

// BuildKeyFileIndex aggregates entries by key file path.
//
// Files with 2+ sessions are marked popular. Sorted by count desc, then path.
//
// Parameters:
//   - entries: All journal entries to aggregate
//
// Returns:
//   - []KeyFileData: Key files sorted by session count descending, then path
func BuildKeyFileIndex(entries []JournalEntry) []KeyFileData {
	grouped := BuildGroupedIndex(
		entries,
		func(e JournalEntry) []string { return e.KeyFiles },
	)
	files := make([]KeyFileData, len(grouped))
	for i, g := range grouped {
		files[i] = KeyFileData{Path: g.Key, Entries: g.Entries, Popular: g.Popular}
	}
	return files
}

// GenerateKeyFilesIndex creates the files/index.md page.
//
// Frequently touched files link to dedicated pages; single-session files
// list entries inline.
//
// Parameters:
//   - keyFiles: Sorted key file data from BuildKeyFileIndex
//
// Returns:
//   - string: Markdown content for files/index.md
func GenerateKeyFilesIndex(keyFiles []KeyFileData) string {
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

	sb.WriteString(desc.TextDesc(text.DescKeyHeadingKeyFiles) + nl + nl)
	sb.WriteString(fmt.Sprintf(
		tpl.TplJournalFileStats+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail)),
	)

	WritePopularAndLongtail(&sb,
		len(popular), desc.TextDesc(text.DescKeyHeadingFrequentlyTouched),
		func(i int) (string, string, int) {
			return "`" + popular[i].Path + "`",
				KeyFileSlug(popular[i].Path),
				len(popular[i].Entries)
		},
		len(longtail), desc.TextDesc(text.DescKeyHeadingSingleSession),
		tpl.TplJournalLongtailCodeEntry,
		func(i int) (string, JournalEntry) {
			return longtail[i].Path, longtail[i].Entries[0]
		},
	)

	return sb.String()
}

// GenerateKeyFilePage creates an individual key file page with sessions
// grouped by month.
//
// Parameters:
//   - kf: Key file data including the path and entries
//
// Returns:
//   - string: Markdown content for the key file page
func GenerateKeyFilePage(kf KeyFileData) string {
	return GenerateGroupedPage(
		fmt.Sprintf(tpl.TplJournalCodePageHeading, kf.Path),
		fmt.Sprintf(tpl.TplJournalFilePageStats, len(kf.Entries)),
		kf.Entries,
	)
}

// BuildTypeIndex aggregates entries by session type.
//
// Sorted by count descending, then name alphabetically.
//
// Parameters:
//   - entries: All journal entries to aggregate
//
// Returns:
//   - []TypeData: Session types sorted by count descending, then name
func BuildTypeIndex(entries []JournalEntry) []TypeData {
	grouped := BuildGroupedIndex(
		entries,
		func(e JournalEntry) []string { return []string{e.Type} },
	)
	types := make([]TypeData, len(grouped))
	for i, g := range grouped {
		types[i] = TypeData{Name: g.Key, Entries: g.Entries}
	}
	return types
}

// GenerateTypesIndex creates the types/index.md page.
//
// Parameters:
//   - sessionTypes: Sorted type data from BuildTypeIndex
//
// Returns:
//   - string: Markdown content for types/index.md
func GenerateTypesIndex(sessionTypes []TypeData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	totalSessions := 0
	for _, st := range sessionTypes {
		totalSessions += len(st.Entries)
	}

	sb.WriteString(desc.TextDesc(text.DescKeyHeadingSessionTypes) + nl + nl)
	sb.WriteString(fmt.Sprintf(
		tpl.TplJournalTypeStats+nl+nl, len(sessionTypes), totalSessions),
	)

	for _, st := range sessionTypes {
		sb.WriteString(FormatSessionLink(st.Name, st.Name, len(st.Entries)))
	}
	sb.WriteString(nl)

	return sb.String()
}

// GenerateTypePage creates an individual session type page with sessions
// grouped by month.
//
// Parameters:
//   - st: Type data including name and entries
//
// Returns:
//   - string: Markdown content for the session type page
func GenerateTypePage(st TypeData) string {
	return GenerateGroupedPage(
		fmt.Sprintf(tpl.TplJournalPageHeading, st.Name),
		fmt.Sprintf(tpl.TplJournalTypePageStats, len(st.Entries), st.Name),
		st.Entries,
	)
}
