//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import "strings"

// WordSet splits text into a set of unique words for O(1) lookup.
//
// Parameters:
//   - text: Input text to split into words
//
// Returns:
//   - map[string]bool: Set of unique words for O(1) membership lookup
func WordSet(text string) map[string]bool {
	fields := strings.Fields(text)
	set := make(map[string]bool, len(fields))
	for _, w := range fields {
		set[w] = true
	}
	return set
}
