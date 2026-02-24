//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"sort"

	"github.com/ActiveMemory/ctx/internal/config"
)

// groupByMonth groups journal entries by their YYYY-MM date prefix,
// preserving insertion order of months.
//
// Parameters:
//   - entries: Journal entries to group (must have Date field set)
//
// Returns:
//   - map[string][]journalEntry: Entries keyed by month string
//   - []string: Month strings in first-seen order
func groupByMonth(
	entries []journalEntry,
) (map[string][]journalEntry, []string) {
	months := make(map[string][]journalEntry)
	var monthOrder []string

	for _, e := range entries {
		if len(e.Date) >= config.JournalMonthPrefixLen {
			month := e.Date[:config.JournalMonthPrefixLen]
			if _, exists := months[month]; !exists {
				monthOrder = append(monthOrder, month)
			}
			months[month] = append(months[month], e)
		}
	}

	return months, monthOrder
}

// buildGroupedIndex aggregates entries by keys extracted via extractKeys,
// marks groups with 2+ sessions as popular, and sorts by count descending
// then alphabetically.
//
// Parameters:
//   - entries: Journal entries to aggregate
//   - extractKeys: Function that returns grouping keys for a given entry
//
// Returns:
//   - []groupedIndex: Sorted groups with popularity flags
func buildGroupedIndex(
	entries []journalEntry, extractKeys func(journalEntry) []string,
) []groupedIndex {
	byKey := make(map[string][]journalEntry)
	for _, e := range entries {
		for _, k := range extractKeys(e) {
			byKey[k] = append(byKey[k], e)
		}
	}

	result := make([]groupedIndex, 0, len(byKey))
	for key, ents := range byKey {
		result = append(result, groupedIndex{
			Key:     key,
			Entries: ents,
			Popular: len(ents) >= config.JournalPopularityThreshold,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if len(result[i].Entries) != len(result[j].Entries) {
			return len(result[i].Entries) > len(result[j].Entries)
		}
		return result[i].Key < result[j].Key
	})

	return result
}
