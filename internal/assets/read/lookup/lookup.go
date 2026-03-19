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
func loadStopWords() map[string]bool {
	raw := TextDesc(text.TextDescKeyStopwords)
	words := strings.Fields(raw)
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}
