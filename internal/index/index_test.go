//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Entry
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
			expected: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "No custom UI - IDE is the interface"},
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
			expected: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "First decision"},
				{Timestamp: "2026-01-27-123456", Date: "2026-01-27", Title: "Second decision"},
			},
		},
		{
			name: "entry with special characters",
			content: `# Decisions

## [2026-01-28-051426] Use tool-agnostic Session type | with pipe

**Status**: Accepted
`,
			expected: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "Use tool-agnostic Session type | with pipe"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseHeaders(tt.content)
			if len(got) != len(tt.expected) {
				t.Errorf("ParseHeaders() got %d entries, want %d", len(got), len(tt.expected))
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

func TestGenerateTable(t *testing.T) {
	tests := []struct {
		name         string
		entries      []Entry
		columnHeader string
		expected     string
	}{
		{
			name:         "empty entries",
			entries:      nil,
			columnHeader: "Decision",
			expected:     "",
		},
		{
			name:         "empty slice",
			entries:      []Entry{},
			columnHeader: "Decision",
			expected:     "",
		},
		{
			name: "single entry",
			entries: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "First decision"},
			},
			columnHeader: "Decision",
			expected: `| Date | Decision |
|------|--------|
| 2026-01-28 | First decision |
`,
		},
		{
			name: "multiple entries",
			entries: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "First"},
				{Timestamp: "2026-01-27-123456", Date: "2026-01-27", Title: "Second"},
			},
			columnHeader: "Decision",
			expected: `| Date | Decision |
|------|--------|
| 2026-01-28 | First |
| 2026-01-27 | Second |
`,
		},
		{
			name: "entry with pipe character",
			entries: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "Use A | B format"},
			},
			columnHeader: "Decision",
			expected: `| Date | Decision |
|------|--------|
| 2026-01-28 | Use A \| B format |
`,
		},
		{
			name: "learning column header",
			entries: []Entry{
				{Timestamp: "2026-01-28-051426", Date: "2026-01-28", Title: "Test entry"},
			},
			columnHeader: "Learning",
			expected: `| Date | Learning |
|------|--------|
| 2026-01-28 | Test entry |
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTable(tt.entries, tt.columnHeader)
			if got != tt.expected {
				t.Errorf("GenerateTable() =\n%q\nwant\n%q", got, tt.expected)
			}
		})
	}
}

func TestUpdateDecisions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantHas []string
		wantNot []string
	}{
		{
			name:    "empty file with header",
			content: "# Decisions\n",
			wantNot: []string{config.IndexStart, config.IndexEnd},
		},
		{
			name: "file with one decision",
			content: `# Decisions

## [2026-01-28-051426] Test decision

**Status**: Accepted
`,
			wantHas: []string{
				config.IndexStart,
				config.IndexEnd,
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
				config.IndexStart,
				config.IndexEnd,
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
				config.IndexStart,
				config.IndexEnd,
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
			got := UpdateDecisions(tt.content)
			for _, want := range tt.wantHas {
				if !strings.Contains(got, want) {
					t.Errorf("UpdateDecisions() result missing %q\nGot:\n%s", want, got)
				}
			}
			for _, notWant := range tt.wantNot {
				if strings.Contains(got, notWant) {
					t.Errorf("UpdateDecisions() result should not contain %q\nGot:\n%s", notWant, got)
				}
			}
		})
	}
}

func TestUpdateDecisions_PreservesContent(t *testing.T) {
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

	got := UpdateDecisions(content)

	if !strings.Contains(got, config.IndexStart) {
		t.Error("Missing INDEX:START marker")
	}
	if !strings.Contains(got, config.IndexEnd) {
		t.Error("Missing INDEX:END marker")
	}

	if !strings.Contains(got, "| 2026-01-28 | First decision |") {
		t.Error("Missing first decision in index")
	}
	if !strings.Contains(got, "| 2026-01-27 | Second decision |") {
		t.Error("Missing second decision in index")
	}

	if !strings.Contains(got, "**Context**: Some context here.") {
		t.Error("Lost content from first decision")
	}
	if !strings.Contains(got, "**Rationale**: Another rationale.") {
		t.Error("Lost content from second decision")
	}
}

func TestUpdateDecisions_Idempotent(t *testing.T) {
	content := `# Decisions

## [2026-01-28-051426] Test decision

**Status**: Accepted
`

	first := UpdateDecisions(content)
	second := UpdateDecisions(first)

	if first != second {
		t.Errorf("UpdateDecisions is not idempotent\nFirst:\n%s\nSecond:\n%s", first, second)
	}
}

func TestUpdateLearnings(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantHas []string
		wantNot []string
	}{
		{
			name:    "empty file with header",
			content: "# Learnings\n",
			wantNot: []string{config.IndexStart, config.IndexEnd},
		},
		{
			name: "file with one learning",
			content: `# Learnings

## [2026-01-28-191951] Required flags now enforced

**Context**: Implemented ctx add learning flags

**Lesson**: Structured entries are more useful

**Application**: Always use all three flags
`,
			wantHas: []string{
				config.IndexStart,
				config.IndexEnd,
				"| Date | Learning |",
				"| 2026-01-28 | Required flags now enforced |",
			},
		},
		{
			name: "multiple learnings",
			content: `# Learnings

## [2026-01-28-191951] First learning

**Context**: Test

**Lesson**: Test

**Application**: Test

---

## [2026-01-27-120000] Second learning

**Context**: Test

**Lesson**: Test

**Application**: Test
`,
			wantHas: []string{
				"| 2026-01-28 | First learning |",
				"| 2026-01-27 | Second learning |",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateLearnings(tt.content)
			for _, want := range tt.wantHas {
				if !strings.Contains(got, want) {
					t.Errorf("UpdateLearnings() result missing %q\nGot:\n%s", want, got)
				}
			}
			for _, notWant := range tt.wantNot {
				if strings.Contains(got, notWant) {
					t.Errorf("UpdateLearnings() result should not contain %q\nGot:\n%s", notWant, got)
				}
			}
		})
	}
}

func TestUpdateLearnings_Idempotent(t *testing.T) {
	content := `# Learnings

## [2026-01-28-191951] Test learning

**Context**: Test

**Lesson**: Test

**Application**: Test
`

	first := UpdateLearnings(content)
	second := UpdateLearnings(first)

	if first != second {
		t.Errorf("UpdateLearnings is not idempotent\nFirst:\n%s\nSecond:\n%s", first, second)
	}
}
