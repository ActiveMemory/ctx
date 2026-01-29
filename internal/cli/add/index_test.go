//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"strings"
	"testing"
)

func TestParseDecisionHeaders(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []DecisionEntry
	}{
		{
			name:     "empty content",
			content:  "",
			expected: nil,
		},
		{
			name:     "no decisions",
			content:  "# Decisions\n\nSome text here.",
			expected: nil,
		},
		{
			name: "single decision",
			content: `# Decisions

## [2026-01-28-051426] No custom UI - IDE is the interface

**Status**: Accepted
`,
			expected: []DecisionEntry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "No custom UI - IDE is the interface"},
			},
		},
		{
			name: "multiple decisions",
			content: `# Decisions

## [2026-01-28-051426] First decision

**Status**: Accepted

---

## [2026-01-27-123456] Second decision

**Status**: Accepted
`,
			expected: []DecisionEntry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "First decision"},
				{Timestamp: "2026-01-27-123456", Date: "2026-01-27", Title: "Second decision"},
			},
		},
		{
			name: "decision with special characters",
			content: `# Decisions

## [2026-01-28-051426] Use tool-agnostic Session type | with pipe

**Status**: Accepted
`,
			expected: []DecisionEntry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "Use tool-agnostic Session type | with pipe"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseDecisionHeaders(tt.content)
			if len(got) != len(tt.expected) {
				t.Errorf("ParseDecisionHeaders() got %d entries, want %d", len(got), len(tt.expected))
				return
			}
			for i, entry := range got {
				if entry.Timestamp != tt.expected[i].Timestamp {
					t.Errorf("entry[%d].Timestamp = %q, want %q", i, entry.Timestamp, tt.expected[i].Timestamp)
				}
				if entry.Date != tt.expected[i].Date {
					t.Errorf("entry[%d].Date = %q, want %q", i, entry.Date, tt.expected[i].Date)
				}
				if entry.Title != tt.expected[i].Title {
					t.Errorf("entry[%d].Title = %q, want %q", i, entry.Title, tt.expected[i].Title)
				}
			}
		})
	}
}

func TestGenerateIndex(t *testing.T) {
	tests := []struct {
		name     string
		entries  []DecisionEntry
		expected string
	}{
		{
			name:     "empty entries",
			entries:  nil,
			expected: "",
		},
		{
			name:     "empty slice",
			entries:  []DecisionEntry{},
			expected: "",
		},
		{
			name: "single entry",
			entries: []DecisionEntry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "First decision"},
			},
			expected: `| Date | Decision |
|------|----------|
| 2026-01-28 | First decision |
`,
		},
		{
			name: "multiple entries",
			entries: []DecisionEntry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "First"},
				{Timestamp: "2026-01-27-123456", Date: "2026-01-27", Title: "Second"},
			},
			expected: `| Date | Decision |
|------|----------|
| 2026-01-28 | First |
| 2026-01-27 | Second |
`,
		},
		{
			name: "entry with pipe character",
			entries: []DecisionEntry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "Use A | B format"},
			},
			expected: `| Date | Decision |
|------|----------|
| 2026-01-28 | Use A \| B format |
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateIndex(tt.entries)
			if got != tt.expected {
				t.Errorf("GenerateIndex() =\n%q\nwant\n%q", got, tt.expected)
			}
		})
	}
}

func TestUpdateIndex(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantHas  []string // strings that should be present
		wantNot  []string // strings that should NOT be present
	}{
		{
			name:    "empty file with header",
			content: "# Decisions\n",
			wantNot: []string{IndexStart, IndexEnd},
		},
		{
			name: "file with one decision",
			content: `# Decisions

## [2026-01-28-051426] Test decision

**Status**: Accepted
`,
			wantHas: []string{
				IndexStart,
				IndexEnd,
				"| Date | Decision |",
				"| 2026-01-28 | Test decision |",
				"## [2026-01-28-051426] Test decision",
			},
		},
		{
			name: "update existing index",
			content: `# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|----------|
| 2026-01-28 | Old entry |
<!-- INDEX:END -->

## [2026-01-28-051426] New decision

**Status**: Accepted
`,
			wantHas: []string{
				IndexStart,
				IndexEnd,
				"| 2026-01-28 | New decision |",
			},
			wantNot: []string{
				"| 2026-01-28 | Old entry |",
			},
		},
		{
			name: "remove index when no decisions",
			content: `# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|----------|
| 2026-01-28 | Old entry |
<!-- INDEX:END -->

Some other content.
`,
			wantNot: []string{
				IndexStart,
				IndexEnd,
				"| Date | Decision |",
			},
			wantHas: []string{
				"# Decisions",
				"Some other content.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateIndex(tt.content)
			for _, want := range tt.wantHas {
				if !strings.Contains(got, want) {
					t.Errorf("UpdateIndex() result missing %q\nGot:\n%s", want, got)
				}
			}
			for _, notWant := range tt.wantNot {
				if strings.Contains(got, notWant) {
					t.Errorf("UpdateIndex() result should not contain %q\nGot:\n%s", notWant, got)
				}
			}
		})
	}
}

func TestUpdateIndex_PreservesContent(t *testing.T) {
	content := `# Decisions

## [2026-01-28-051426] First decision

**Status**: Accepted

**Context**: Some context here.

**Decision**: The decision text.

**Rationale**: Why we did it.

**Consequences**: What happens next.

---

## [2026-01-27-123456] Second decision

**Status**: Accepted

**Context**: Another context.

**Decision**: Another decision.

**Rationale**: Another rationale.

**Consequences**: More consequences.
`

	got := UpdateIndex(content)

	// Index should be present
	if !strings.Contains(got, IndexStart) {
		t.Error("Missing INDEX:START marker")
	}
	if !strings.Contains(got, IndexEnd) {
		t.Error("Missing INDEX:END marker")
	}

	// Both entries should be in index
	if !strings.Contains(got, "| 2026-01-28 | First decision |") {
		t.Error("Missing first decision in index")
	}
	if !strings.Contains(got, "| 2026-01-27 | Second decision |") {
		t.Error("Missing second decision in index")
	}

	// Full content should be preserved
	if !strings.Contains(got, "**Context**: Some context here.") {
		t.Error("Lost content from first decision")
	}
	if !strings.Contains(got, "**Rationale**: Another rationale.") {
		t.Error("Lost content from second decision")
	}
}

func TestUpdateIndex_Idempotent(t *testing.T) {
	content := `# Decisions

## [2026-01-28-051426] Test decision

**Status**: Accepted
`

	// Apply once
	first := UpdateIndex(content)

	// Apply again
	second := UpdateIndex(first)

	// Should be identical
	if first != second {
		t.Errorf("UpdateIndex is not idempotent\nFirst:\n%s\nSecond:\n%s", first, second)
	}
}
