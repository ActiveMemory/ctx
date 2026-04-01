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

// loadStopWords parses the stopwords text entry into a lookup map.
//
// Returns:
//   - map[string]bool: Set of lowercase stop words
//     keyed for O(1) membership checks
func loadStopWords() map[string]bool {
	raw := TextDesc(text.DescKeyStopwords)
	words := strings.Fields(raw)
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}
