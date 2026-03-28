//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// TextDesc resolves a text description key to its short value.
//
// Parameters:
//   - name: Text key to look up (e.g., "check-persistence.fallback")
//
// Returns:
//   - string: Short description text, or empty string if not found
func TextDesc(name string) string {
	entry, ok := TextMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}

// StopWords returns the default set of stop words for keyword extraction.
//
// Returns:
//   - map[string]bool: Set of lowercase stop words
func StopWords() map[string]bool {
	return stopWordsMap
}

// loadStopWords parses the stopwords text entry into a lookup map.
//
// Returns:
//   - map[string]bool: Set of lowercase stop words keyed for O(1) membership checks
func loadStopWords() map[string]bool {
	raw := TextDesc(text.DescKeyStopwords)
	words := strings.Fields(raw)
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}
