//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"sort"

	"github.com/ActiveMemory/ctx/internal/config/journal"
)

// GroupByMonth groups journal entries by their YYYY-MM date prefix,
// preserving insertion order of months.
//
// Parameters:
//   - entries: Journal entries to group (must have Date field set)
//
// Returns:
//   - map[string][]JournalEntry: Entries keyed by month string
//   - []string: Month strings in first-seen order
func GroupByMonth(
	entries []JournalEntry,
) (map[string][]JournalEntry, []string) {
	months := make(map[string][]JournalEntry)
	var monthOrder []string

	for _, e := range entries {
		if len(e.Date) >= journal.MonthPrefixLen {
			month := e.Date[:journal.MonthPrefixLen]
			if _, exists := months[month]; !exists {
				monthOrder = append(monthOrder, month)
			}
			months[month] = append(months[month], e)
		}
	}

	return months, monthOrder
}

// BuildGroupedIndex aggregates entries by keys extracted via extractKeys,
// marks groups with 2+ sessions as popular, and sorts by count descending
// then alphabetically.
//
// Parameters:
//   - entries: Journal entries to aggregate
//   - extractKeys: Function that returns grouping keys for a given entry
//
// Returns:
//   - []GroupedIndex: Sorted groups with popularity flags
func BuildGroupedIndex(
	entries []JournalEntry, extractKeys func(JournalEntry) []string,
) []GroupedIndex {
	byKey := make(map[string][]JournalEntry)
	for _, e := range entries {
		for _, k := range extractKeys(e) {
			byKey[k] = append(byKey[k], e)
		}
	}

	result := make([]GroupedIndex, 0, len(byKey))
	for key, ents := range byKey {
		result = append(result, GroupedIndex{
			Key:     key,
			Entries: ents,
			Popular: len(ents) >= journal.PopularityThreshold,
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
