//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package section

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/format"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/group"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
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
func BuildTopicIndex(entries []entity.JournalEntry) []entity.TopicData {
	grouped := group.GroupedIndex(
		entries,
		func(e entity.JournalEntry) []string { return e.Topics },
	)
	topics := make([]entity.TopicData, len(grouped))
	for i, g := range grouped {
		topics[i] = entity.TopicData{
			Name: g.Key, Entries: g.Entries,
			Popular: g.Popular,
		}
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
func GenerateTopicsIndex(topics []entity.TopicData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	var popular, longtail []entity.TopicData
	for _, t := range topics {
		if t.Popular {
			popular = append(popular, t)
		} else {
			longtail = append(longtail, t)
		}
	}

	sb.WriteString(desc.Text(text.DescKeyHeadingTopics) + nl + nl)
	io.SafeFprintf(&sb,
		tpl.JournalTopicStats+nl+nl,
		len(topics), session.CountUnique(topics), len(popular), len(longtail))

	WritePopularAndLongtail(&sb,
		len(popular), desc.Text(text.DescKeyHeadingPopularTopics),
		func(i int) (string, string, int) {
			return popular[i].Name, popular[i].Name, len(popular[i].Entries)
		},
		len(longtail), desc.Text(text.DescKeyHeadingLongtailTopics),
		tpl.JournalLongtailEntry,
		func(i int) (string, entity.JournalEntry) {
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
func GenerateTopicPage(topic entity.TopicData) string {
	return GenerateGroupedPage(
		fmt.Sprintf(tpl.JournalPageHeading, topic.Name),
		fmt.Sprintf(tpl.JournalTopicPageStats, len(topic.Entries)),
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
func BuildKeyFileIndex(entries []entity.JournalEntry) []entity.KeyFileData {
	grouped := group.GroupedIndex(
		entries,
		func(e entity.JournalEntry) []string { return e.KeyFiles },
	)
	files := make([]entity.KeyFileData, len(grouped))
	for i, g := range grouped {
		files[i] = entity.KeyFileData{
			Path: g.Key, Entries: g.Entries,
			Popular: g.Popular,
		}
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
func GenerateKeyFilesIndex(keyFiles []entity.KeyFileData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	var popular, longtail []entity.KeyFileData
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

	sb.WriteString(desc.Text(text.DescKeyHeadingKeyFiles) + nl + nl)
	io.SafeFprintf(&sb,
		tpl.JournalFileStats+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail))

	WritePopularAndLongtail(&sb,
		len(popular), desc.Text(text.DescKeyHeadingFrequentlyTouched),
		func(i int) (string, string, int) {
			return token.Backtick + popular[i].Path + token.Backtick,
				format.KeyFileSlug(popular[i].Path),
				len(popular[i].Entries)
		},
		len(longtail), desc.Text(text.DescKeyHeadingSingleSession),
		tpl.JournalLongtailCodeEntry,
		func(i int) (string, entity.JournalEntry) {
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
func GenerateKeyFilePage(kf entity.KeyFileData) string {
	return GenerateGroupedPage(
		fmt.Sprintf(tpl.JournalCodePageHeading, kf.Path),
		fmt.Sprintf(tpl.JournalFilePageStats, len(kf.Entries)),
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
func BuildTypeIndex(entries []entity.JournalEntry) []entity.TypeData {
	grouped := group.GroupedIndex(
		entries,
		func(e entity.JournalEntry) []string { return []string{e.Type} },
	)
	types := make([]entity.TypeData, len(grouped))
	for i, g := range grouped {
		types[i] = entity.TypeData{Name: g.Key, Entries: g.Entries}
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
func GenerateTypesIndex(sessionTypes []entity.TypeData) string {
	var sb strings.Builder
	nl := token.NewlineLF

	totalSessions := 0
	for _, st := range sessionTypes {
		totalSessions += len(st.Entries)
	}

	sb.WriteString(desc.Text(text.DescKeyHeadingSessionTypes) + nl + nl)
	io.SafeFprintf(&sb,
		tpl.JournalTypeStats+nl+nl, len(sessionTypes), totalSessions)

	for _, st := range sessionTypes {
		sb.WriteString(format.SessionLink(st.Name, st.Name, len(st.Entries)))
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
func GenerateTypePage(st entity.TypeData) string {
	return GenerateGroupedPage(
		fmt.Sprintf(tpl.JournalPageHeading, st.Name),
		fmt.Sprintf(tpl.JournalTypePageStats, len(st.Entries), st.Name),
		st.Entries,
	)
}
