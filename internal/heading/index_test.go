//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heading

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/entity"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []entity.IndexEntry
	}{
		{
			name:     "empty content",
			content:  "",
			expected: nil,
		},
		{
			name:     "no entries",
			content:  "# Decisions\n\nSome text here.",
			expected: nil,
		},
		{
			name: "single entry",
			content: `# Decisions

## [2026-01-28-051426] No custom UI - IDE is the interface

**Status**: Accepted
`,
			expected: []entity.IndexEntry{
				{
					Timestamp: "2026-01-28-051426",
					Date:      "2026-01-28",
					Title:     "No custom UI - IDE is the interface",
				},
			},
		},
		{
			name: "multiple entries",
			content: `# Decisions

## [2026-01-28-051426] First decision

**Status**: Accepted

---

## [2026-01-27-123456] Second decision

**Status**: Accepted
`,
			expected: []entity.IndexEntry{
				{
					Timestamp: "2026-01-28-051426",
					Date:      "2026-01-28",
					Title:     "First decision",
				},
				{
					Timestamp: "2026-01-27-123456",
					Date:      "2026-01-27",
					Title:     "Second decision",
				},
			},
		},
		{
			name: "entry with special characters",
			content: `# Decisions

## [2026-01-28-051426] Use tool-agnostic Session type | with pipe

**Status**: Accepted
`,
			expected: []entity.IndexEntry{
				{
					Timestamp: "2026-01-28-051426",
					Date:      "2026-01-28",
					Title:     "Use tool-agnostic Session type | with pipe",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseHeaders(tt.content)
			if len(got) != len(tt.expected) {
				t.Errorf(
					"ParseHeaders() got %d entries, want %d",
					len(got), len(tt.expected),
				)
				return
			}
			for i, entry := range got {
				if entry.Timestamp != tt.expected[i].Timestamp {
					t.Errorf(
						"entry[%d].Timestamp = %q, want %q",
						i, entry.Timestamp,
						tt.expected[i].Timestamp,
					)
				}
				if entry.Date != tt.expected[i].Date {
					t.Errorf(
						"entry[%d].Date = %q, want %q",
						i, entry.Date, tt.expected[i].Date,
					)
				}
				if entry.Title != tt.expected[i].Title {
					t.Errorf(
						"entry[%d].Title = %q, want %q",
						i, entry.Title, tt.expected[i].Title,
					)
				}
			}
		})
	}
}
