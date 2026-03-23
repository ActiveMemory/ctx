//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/hook"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// FormatStaleEntries builds the display text for stale files.
//
// Groups entries by review URL. Entries without a URL get the generic
// "touch to mark as reviewed" footer. Entries with URLs get a
// "Review against: <url>" line after their group.
//
// Parameters:
//   - entries: stale file entries to format
//
// Returns:
//   - string: formatted multi-line text for display
func FormatStaleEntries(entries []hook.StaleEntry) string {
	byURL := make(map[string][]hook.StaleEntry)
	var noURL []hook.StaleEntry
	var urlOrder []string

	for _, e := range entries {
		if e.ReviewURL == "" {
			noURL = append(noURL, e)
			continue
		}
		if _, seen := byURL[e.ReviewURL]; !seen {
			urlOrder = append(urlOrder, e.ReviewURL)
		}
		byURL[e.ReviewURL] = append(byURL[e.ReviewURL], e)
	}

	var b strings.Builder

	for _, url := range urlOrder {
		group := byURL[url]
		for _, e := range group {
			_, writeErr := fmt.Fprintf(&b, desc.Text(text.DescKeyFreshnessFileEntry),
				e.Path, e.Days, e.Desc)
			if writeErr != nil {
				return ""
			}
			b.WriteString(token.NewlineLF)
		}
		_, writeErr := fmt.Fprintf(&b, desc.Text(text.DescKeyFreshnessReviewURL), url)
		if writeErr != nil {
			return ""
		}
		b.WriteString(token.NewlineLF)
	}

	for _, e := range noURL {
		_, writeErr := fmt.Fprintf(&b, desc.Text(text.DescKeyFreshnessFileEntry),
			e.Path, e.Days, e.Desc)
		if writeErr != nil {
			return ""
		}
		b.WriteString(token.NewlineLF)
	}

	b.WriteString(desc.Text(text.DescKeyFreshnessTouchHint))

	return b.String()
}
