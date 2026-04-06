//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/config/regex"
)

// Extract returns all unique tags from an entry string, sorted
// alphabetically. For blob entries, only the label portion is scanned.
//
// Parameters:
//   - entry: Scratchpad entry string
//
// Returns:
//   - []string: Unique tag names without the # prefix, or nil
func Extract(entry string) []string {
	text := ScanText(entry)
	matches := regex.Tag.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	var tags []string
	for _, m := range matches {
		name := m[1]
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		tags = append(tags, name)
	}
	sort.Strings(tags)
	return tags
}

// Has returns true if the entry contains the given tag.
//
// Parameters:
//   - entry: Scratchpad entry string
//   - tagName: Tag name without the # prefix
//
// Returns:
//   - bool: True if the tag is present
func Has(entry, tagName string) bool {
	text := ScanText(entry)
	matches := regex.Tag.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		if m[1] == tagName {
			return true
		}
	}
	return false
}

// Match returns true if the entry satisfies the filter.
// A filter prefixed with "~" negates the match: "~later" matches
// entries that do NOT contain #later.
//
// Parameters:
//   - entry: Scratchpad entry string
//   - filter: Tag name, optionally prefixed with "~" for negation
//
// Returns:
//   - bool: True if the entry matches the filter
func Match(entry, filter string) bool {
	if strings.HasPrefix(filter, pad.TagNegate) {
		return !Has(entry, filter[len(pad.TagNegate):])
	}
	return Has(entry, filter)
}

// MatchAll returns true if the entry satisfies all filters (AND logic).
//
// Parameters:
//   - entry: Scratchpad entry string
//   - filters: Tag filters, each optionally prefixed with "!"
//
// Returns:
//   - bool: True if all filters match
func MatchAll(entry string, filters []string) bool {
	for _, f := range filters {
		if !Match(entry, f) {
			return false
		}
	}
	return true
}
