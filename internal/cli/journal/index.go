//   /    Context:                     https://ctx.ist
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

// buildTopicIndex aggregates entries by topic and returns sorted topic data.
//
// Topics with 2+ sessions are marked popular. Sorted by count desc, then alpha.
//
// Parameters:
//   - entries: All journal entries to aggregate
//
// Returns:
//   - []topicData: Topics sorted by session count descending, then name
func buildTopicIndex(entries []journalEntry) []topicData {
	grouped := buildGroupedIndex(
		entries,
		func(e journalEntry) []string { return e.Topics },
	)
	topics := make([]topicData, len(grouped))
	for i, g := range grouped {
		topics[i] = topicData{Name: g.Key, Entries: g.Entries, Popular: g.Popular}
	}
	return topics
}

// generateTopicsIndex creates the topics/index.md page.
//
// Popular topics link to dedicated pages; long-tail topics list entries inline.
//
// Parameters:
//   - topics: Sorted topic data from buildTopicIndex
//
// Returns:
//   - string: Markdown content for topics/index.md
func generateTopicsIndex(topics []topicData) string {
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

	sb.WriteString(config.JournalHeadingTopics + nl + nl)
	sb.WriteString(fmt.Sprintf(
		config.TplJournalTopicStats+nl+nl,
		len(topics), countUniqueSessions(topics), len(popular), len(longtail)))

	writePopularAndLongtail(&sb,
		len(popular), config.JournalHeadingPopularTopics,
		func(i int) (string, string, int) {
			return popular[i].Name, popular[i].Name, len(popular[i].Entries)
		},
		len(longtail), config.JournalHeadingLongtailTopics,
		config.TplJournalLongtailEntry,
		func(i int) (string, journalEntry) {
			return longtail[i].Name, longtail[i].Entries[0]
		},
	)

	return sb.String()
}

// generateTopicPage creates an individual topic page with sessions grouped
// by month.
//
// Parameters:
//   - topic: Topic data including name and entries
//
// Returns:
//   - string: Markdown content for the topic page
func generateTopicPage(topic topicData) string {
	return generateGroupedPage(
		fmt.Sprintf(config.TplJournalPageHeading, topic.Name),
		fmt.Sprintf(config.TplJournalTopicPageStats, len(topic.Entries)),
		topic.Entries,
	)
}

// buildKeyFileIndex aggregates entries by key file path.
//
// Files with 2+ sessions are marked popular. Sorted by count desc, then path.
//
// Parameters:
//   - entries: All journal entries to aggregate
//
// Returns:
//   - []keyFileData: Key files sorted by session count descending, then path
func buildKeyFileIndex(entries []journalEntry) []keyFileData {
	grouped := buildGroupedIndex(
		entries,
		func(e journalEntry) []string { return e.KeyFiles },
	)
	files := make([]keyFileData, len(grouped))
	for i, g := range grouped {
		files[i] = keyFileData{Path: g.Key, Entries: g.Entries, Popular: g.Popular}
	}
	return files
}

// generateKeyFilesIndex creates the files/index.md page.
//
// Frequently touched files link to dedicated pages; single-session files
// list entries inline.
//
// Parameters:
//   - keyFiles: Sorted key file data from buildKeyFileIndex
//
// Returns:
//   - string: Markdown content for files/index.md
func generateKeyFilesIndex(keyFiles []keyFileData) string {
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

	sb.WriteString(config.JournalHeadingKeyFiles + nl + nl)
	sb.WriteString(fmt.Sprintf(
		config.TplJournalFileStats+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail)),
	)

	writePopularAndLongtail(&sb,
		len(popular), config.JournalHeadingFrequentlyTouched,
		func(i int) (string, string, int) {
			return "`" + popular[i].Path + "`",
				keyFileSlug(popular[i].Path),
				len(popular[i].Entries)
		},
		len(longtail), config.JournalHeadingSingleSession,
		config.TplJournalLongtailCodeEntry,
		func(i int) (string, journalEntry) {
			return longtail[i].Path, longtail[i].Entries[0]
		},
	)

	return sb.String()
}

// generateKeyFilePage creates an individual key file page with sessions
// grouped by month.
//
// Parameters:
//   - kf: Key file data including the path and entries
//
// Returns:
//   - string: Markdown content for the key file page
func generateKeyFilePage(kf keyFileData) string {
	return generateGroupedPage(
		fmt.Sprintf(config.TplJournalCodePageHeading, kf.Path),
		fmt.Sprintf(config.TplJournalFilePageStats, len(kf.Entries)),
		kf.Entries,
	)
}

// buildTypeIndex aggregates entries by session type.
//
// Sorted by count descending, then name alphabetically.
//
// Parameters:
//   - entries: All journal entries to aggregate
//
// Returns:
//   - []typeData: Session types sorted by count descending, then name
func buildTypeIndex(entries []journalEntry) []typeData {
	grouped := buildGroupedIndex(
		entries,
		func(e journalEntry) []string { return []string{e.Type} },
	)
	types := make([]typeData, len(grouped))
	for i, g := range grouped {
		types[i] = typeData{Name: g.Key, Entries: g.Entries}
	}
	return types
}

// generateTypesIndex creates the types/index.md page.
//
// Parameters:
//   - sessionTypes: Sorted type data from buildTypeIndex
//
// Returns:
//   - string: Markdown content for types/index.md
func generateTypesIndex(sessionTypes []typeData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	totalSessions := 0
	for _, st := range sessionTypes {
		totalSessions += len(st.Entries)
	}

	sb.WriteString(config.JournalHeadingSessionTypes + nl + nl)
	sb.WriteString(fmt.Sprintf(
		config.TplJournalTypeStats+nl+nl, len(sessionTypes), totalSessions),
	)

	for _, st := range sessionTypes {
		sb.WriteString(formatSessionLink(st.Name, st.Name, len(st.Entries)))
	}
	sb.WriteString(nl)

	return sb.String()
}

// generateTypePage creates an individual session type page with sessions
// grouped by month.
//
// Parameters:
//   - st: Type data including name and entries
//
// Returns:
//   - string: Markdown content for the session type page
func generateTypePage(st typeData) string {
	return generateGroupedPage(
		fmt.Sprintf(config.TplJournalPageHeading, st.Name),
		fmt.Sprintf(config.TplJournalTypePageStats, len(st.Entries), st.Name),
		st.Entries,
	)
}
