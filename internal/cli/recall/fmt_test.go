//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// stubDuration implements the interface{ Minutes() float64 } used by formatDuration.
type stubDuration struct{ mins float64 }

func (d stubDuration) Minutes() float64 { return d.mins }

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		mins float64
		want string
	}{
		{"zero", 0, "<1m"},
		{"sub-minute", 0.5, "<1m"},
		{"one minute", 1, "1m"},
		{"several minutes", 25, "25m"},
		{"fifty-nine minutes", 59, "59m"},
		{"exactly one hour", 60, "1h"},
		{"one hour thirty", 90, "1h30m"},
		{"two hours", 120, "2h"},
		{"two hours fifteen", 135, "2h15m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(stubDuration{tt.mins})
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.mins, got, tt.want)
			}
		})
	}
}

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		name   string
		tokens int
		want   string
	}{
		{"zero", 0, "0"},
		{"small", 500, "500"},
		{"below-K", 999, "999"},
		{"exactly-1K", 1000, "1.0K"},
		{"mid-K", 1500, "1.5K"},
		{"large-K", 50000, "50.0K"},
		{"below-M", 999999, "1000.0K"},
		{"exactly-1M", 1000000, "1.0M"},
		{"mid-M", 2300000, "2.3M"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTokens(tt.tokens)
			if got != tt.want {
				t.Errorf("formatTokens(%d) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestFenceForContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"no backticks", "hello world", "```"},
		{"single backtick", "use `code` here", "```"},
		{"triple backticks", "```go\nfmt.Println()\n```", "````"},
		{"quad backticks", "````\ncode\n````", "`````"},
		{"nested fences", "text\n```\ninner\n```\nmore", "````"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fenceForContent(tt.content)
			if got != tt.want {
				t.Errorf("fenceForContent(%q) = %q, want %q", tt.content, got, tt.want)
			}
		})
	}
}

func TestStripLineNumbers(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"no line numbers", "hello\nworld", "hello\nworld"},
		{"with line numbers", "  1→hello\n  2→world", "hello\nworld"},
		{"mixed", "  1→first\nplain\n  3→third", "first\nplain\nthird"},
		{"large numbers", "  100→line hundred\n  101→next", "line hundred\nnext"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripLineNumbers(tt.content)
			if got != tt.want {
				t.Errorf("stripLineNumbers() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractSystemReminders(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantClean     string
		wantReminders int
	}{
		{
			name:          "no reminders",
			content:       "plain text content",
			wantClean:     "plain text content",
			wantReminders: 0,
		},
		{
			name:          "single reminder",
			content:       "before <system-reminder>reminder text</system-reminder> after",
			wantClean:     "before  after",
			wantReminders: 1,
		},
		{
			name:          "multiple reminders",
			content:       "<system-reminder>first</system-reminder> middle <system-reminder>second</system-reminder>",
			wantClean:     " middle ",
			wantReminders: 2,
		},
		{
			name:          "multiline reminder",
			content:       "text\n<system-reminder>\nmultiline\nreminder\n</system-reminder>\nmore",
			wantClean:     "text\n\nmore",
			wantReminders: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClean, gotReminders := extractSystemReminders(tt.content)
			if gotClean != tt.wantClean {
				t.Errorf("cleaned = %q, want %q", gotClean, tt.wantClean)
			}
			if len(gotReminders) != tt.wantReminders {
				t.Errorf("got %d reminders, want %d", len(gotReminders), tt.wantReminders)
			}
		})
	}
}

func TestNormalizeCodeFences(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "already separated",
			content: "text\n\n```\ncode\n```\n\nmore",
			want:    "text\n\n```\ncode\n```\n\nmore",
		},
		{
			name:    "inline open",
			content: "text ```\ncode\n```",
			want:    "text\n\n```\ncode\n```",
		},
		{
			name:    "close followed by text",
			content: "```\ncode\n``` more text",
			want:    "```\ncode\n```\n\nmore text",
		},
		{
			name:    "both inline",
			content: "before ```\ncode\n``` after",
			want:    "before\n\n```\ncode\n```\n\nafter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeCodeFences(tt.content)
			if got != tt.want {
				t.Errorf("normalizeCodeFences() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestFormatToolUse(t *testing.T) {
	tests := []struct {
		name string
		tool parser.ToolUse
		want string
	}{
		{
			name: "Read tool",
			tool: parser.ToolUse{Name: "Read", Input: `{"file_path":"/tmp/test.go"}`},
			want: "Read: /tmp/test.go",
		},
		{
			name: "Bash tool short",
			tool: parser.ToolUse{Name: "Bash", Input: `{"command":"ls -la"}`},
			want: "Bash: ls -la",
		},
		{
			name: "Bash tool truncated",
			tool: parser.ToolUse{
				Name:  "Bash",
				Input: `{"command":"` + strings.Repeat("x", 150) + `"}`,
			},
			want: "Bash: " + strings.Repeat("x", 100) + "...",
		},
		{
			name: "Grep tool",
			tool: parser.ToolUse{Name: "Grep", Input: `{"pattern":"TODO"}`},
			want: "Grep: TODO",
		},
		{
			name: "unknown tool",
			tool: parser.ToolUse{Name: "CustomTool", Input: `{"anything":"value"}`},
			want: "CustomTool",
		},
		{
			name: "invalid JSON",
			tool: parser.ToolUse{Name: "Read", Input: `not json`},
			want: "Read",
		},
		{
			name: "missing key",
			tool: parser.ToolUse{Name: "Read", Input: `{"other":"value"}`},
			want: "Read",
		},
		{
			name: "Write tool",
			tool: parser.ToolUse{Name: "Write", Input: `{"file_path":"/tmp/out.txt"}`},
			want: "Write: /tmp/out.txt",
		},
		{
			name: "WebSearch tool",
			tool: parser.ToolUse{Name: "WebSearch", Input: `{"query":"golang testing"}`},
			want: "WebSearch: golang testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatToolUse(tt.tool)
			if got != tt.want {
				t.Errorf("formatToolUse() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatPartNavigation(t *testing.T) {
	tests := []struct {
		name       string
		part       int
		totalParts int
		baseName   string
		wantPrev   bool
		wantNext   bool
		wantPartOf string
	}{
		{
			name:       "first of 3",
			part:       1,
			totalParts: 3,
			baseName:   "session",
			wantPrev:   false,
			wantNext:   true,
			wantPartOf: "**Part 1 of 3**",
		},
		{
			name:       "middle of 3",
			part:       2,
			totalParts: 3,
			baseName:   "session",
			wantPrev:   true,
			wantNext:   true,
			wantPartOf: "**Part 2 of 3**",
		},
		{
			name:       "last of 3",
			part:       3,
			totalParts: 3,
			baseName:   "session",
			wantPrev:   true,
			wantNext:   false,
			wantPartOf: "**Part 3 of 3**",
		},
		{
			name:       "part 2 of 2 prev links to base",
			part:       2,
			totalParts: 2,
			baseName:   "session",
			wantPrev:   true,
			wantNext:   false,
			wantPartOf: "**Part 2 of 2**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPartNavigation(tt.part, tt.totalParts, tt.baseName)
			if !strings.Contains(got, tt.wantPartOf) {
				t.Errorf("missing part indicator %q in:\n%s", tt.wantPartOf, got)
			}
			hasPrev := strings.Contains(got, "Previous")
			if hasPrev != tt.wantPrev {
				t.Errorf("hasPrev = %v, want %v", hasPrev, tt.wantPrev)
			}
			hasNext := strings.Contains(got, "Next")
			if hasNext != tt.wantNext {
				t.Errorf("hasNext = %v, want %v", hasNext, tt.wantNext)
			}
			// Part 2 of 2: prev should link to base.md, not p1
			if tt.part == 2 && tt.totalParts == 2 {
				if !strings.Contains(got, tt.baseName+".md") {
					t.Errorf("part 2 of 2 should link prev to %s.md, got:\n%s", tt.baseName, got)
				}
			}
			// Part 3 of 3: prev should link to p2
			if tt.part == 3 && tt.totalParts == 3 {
				if !strings.Contains(got, tt.baseName+"-p2.md") {
					t.Errorf("part 3 of 3 should link prev to %s-p2.md, got:\n%s", tt.baseName, got)
				}
			}
		})
	}
}

// --- formatJournalEntryPart tests ---

func TestFormatJournalEntryPart_SinglePart(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:             "abc12345-session-id",
		Slug:           "test-slug",
		Tool:           "claude-code",
		Project:        "myproject",
		StartTime:      time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		EndTime:        time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
		Duration:       30 * time.Minute,
		TurnCount:      1,
		TotalTokens:    15000,
		TotalTokensIn:  10000,
		TotalTokensOut: 5000,
		Messages: []parser.Message{
			{Role: "user", Text: "Hello", Timestamp: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)},
			{Role: "assistant", Text: "Hi there!", Timestamp: time.Date(2026, 1, 15, 10, 30, 5, 0, time.UTC)},
		},
	}

	got := formatJournalEntryPart(s, s.Messages, 0, 1, 1, "2026-01-15-test-slug-abc12345", "")

	// Verify YAML frontmatter
	if !strings.Contains(got, "---\ndate: \"2026-01-15\"") {
		t.Error("missing YAML frontmatter with date")
	}
	if !strings.Contains(got, "time: \"10:30:00\"") {
		t.Error("missing time in frontmatter")
	}
	if !strings.Contains(got, "project: myproject") {
		t.Error("missing project in frontmatter")
	}
	// Verify slug in heading
	if !strings.Contains(got, "# test-slug") {
		t.Error("missing slug in heading")
	}
	// Verify collapsible metadata table
	if !strings.Contains(got, "<details>") {
		t.Error("missing details block")
	}
	if !strings.Contains(got, "<td><strong>ID</strong></td>") {
		t.Error("missing ID in metadata table")
	}
	if !strings.Contains(got, "<td><strong>Date</strong></td>") {
		t.Error("missing Date in metadata table")
	}
	// Verify token stats in metadata table
	if !strings.Contains(got, "15.0K") {
		t.Error("missing total tokens")
	}
	if !strings.Contains(got, "<td><strong>Tokens</strong></td>") {
		t.Error("missing Tokens in metadata table")
	}
	// Verify conversation content
	if !strings.Contains(got, "Hello") {
		t.Error("missing user message text")
	}
	if !strings.Contains(got, "Hi there!") {
		t.Error("missing assistant message text")
	}
	// NO part navigation for single part
	if strings.Contains(got, "Previous") || strings.Contains(got, "Next") {
		t.Error("single part should have no navigation links")
	}
}

func TestFormatJournalEntryPart_MultiPart(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:             "multi-session-id-12345678",
		Slug:           "multi-part-session",
		Tool:           "claude-code",
		Project:        "proj",
		StartTime:      time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC),
		EndTime:        time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC),
		Duration:       60 * time.Minute,
		TurnCount:      3,
		TotalTokens:    5000,
		TotalTokensIn:  3000,
		TotalTokensOut: 2000,
		Messages: []parser.Message{
			{Role: "user", Text: "msg1", Timestamp: time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC)},
			{Role: "assistant", Text: "resp1", Timestamp: time.Date(2026, 2, 1, 9, 1, 0, 0, time.UTC)},
			{Role: "user", Text: "msg2", Timestamp: time.Date(2026, 2, 1, 9, 5, 0, 0, time.UTC)},
		},
	}

	baseName := "2026-02-01-multi-part-session-multi-se"

	// Part 1 of 3: has metadata + nav
	part1 := formatJournalEntryPart(s, s.Messages[:2], 0, 1, 3, baseName, "")
	if !strings.Contains(part1, "<details>") {
		t.Error("part 1 should have details metadata")
	}
	if !strings.Contains(part1, "<td><strong>ID</strong></td>") {
		t.Error("part 1 should have ID in metadata table")
	}
	if !strings.Contains(part1, "**Part 1 of 3**") {
		t.Error("part 1 should have part indicator")
	}
	if !strings.Contains(part1, "Next") {
		t.Error("part 1 should have next link")
	}

	// Part 2 of 3: no metadata, has nav
	part2 := formatJournalEntryPart(s, s.Messages[2:], 2, 2, 3, baseName, "")
	if strings.Contains(part2, "<details>") {
		t.Error("part 2 should NOT have HTML details metadata")
	}
	if !strings.Contains(part2, "**Part 2 of 3**") {
		t.Error("part 2 should have part indicator")
	}
	if !strings.Contains(part2, "Previous") {
		t.Error("part 2 should have prev link")
	}
	if !strings.Contains(part2, "Next") {
		t.Error("part 2 should have next link")
	}
	// Part 2 should have "continued from part 1"
	if !strings.Contains(part2, "continued from part 1") {
		t.Error("part 2 should indicate continuation")
	}
}

func TestFormatJournalEntryPart_WithToolUse(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:        "tool-session-id-1234",
		Slug:      "tool-session",
		Tool:      "claude-code",
		Project:   "proj",
		StartTime: time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 3, 1, 8, 5, 0, 0, time.UTC),
		Duration:  5 * time.Minute,
		TurnCount: 1,
		Messages: []parser.Message{
			{
				Role:      "user",
				Text:      "Read a file",
				Timestamp: time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC),
			},
			{
				Role:      "assistant",
				Timestamp: time.Date(2026, 3, 1, 8, 0, 5, 0, time.UTC),
				ToolUses: []parser.ToolUse{
					{ID: "t1", Name: "Read", Input: `{"file_path":"/tmp/test.go"}`},
				},
			},
			{
				Role:      "user",
				Timestamp: time.Date(2026, 3, 1, 8, 0, 6, 0, time.UTC),
				ToolResults: []parser.ToolResult{
					{ToolUseID: "t1", Content: "package main\nfunc main() {}"},
				},
			},
			{
				Role:      "user",
				Timestamp: time.Date(2026, 3, 1, 8, 0, 7, 0, time.UTC),
				ToolResults: []parser.ToolResult{
					{ToolUseID: "t2", Content: "error occurred", IsError: true},
				},
			},
			{
				Role:      "user",
				Timestamp: time.Date(2026, 3, 1, 8, 0, 8, 0, time.UTC),
				ToolResults: []parser.ToolResult{
					{
						ToolUseID: "t3",
						Content:   strings.Repeat("line\n", 15), // >10 lines
					},
				},
			},
		},
	}

	got := formatJournalEntryPart(s, s.Messages, 0, 1, 1, "tool-session", "")

	// Verify formatted tool use
	if !strings.Contains(got, "Read: /tmp/test.go") {
		t.Error("missing formatted tool use")
	}
	// Verify tool result in code fence
	if !strings.Contains(got, "package main") {
		t.Error("missing tool result content")
	}
	// Verify error marker
	if !strings.Contains(got, "Error") {
		t.Error("missing error marker for IsError result")
	}
	// Verify collapsible details for long output uses <pre> (not code fences)
	if !strings.Contains(got, "<details>") {
		t.Error("long output (>10 lines) should use <details>")
	}
	if !strings.Contains(got, "</details>") {
		t.Error("long output should have closing </details>")
	}
	if !strings.Contains(got, "<pre>") {
		t.Error("collapsed output should use <pre> tag, not code fences")
	}
}

func TestFormatJournalFilename_WithSlugOverride(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:        "abc12345-full-session-uuid",
		Slug:      "random-claude-slug",
		StartTime: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	// Without override: uses s.Slug.
	got := formatJournalFilename(s, "")
	if got != "2026-01-15-random-claude-slug-abc12345.md" {
		t.Errorf("without override = %q", got)
	}

	// With override: uses the override slug.
	got = formatJournalFilename(s, "fix-auth-bug")
	if got != "2026-01-15-fix-auth-bug-abc12345.md" {
		t.Errorf("with override = %q", got)
	}
}

func TestFormatJournalEntryPart_SessionIDInFrontmatter(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:        "abc12345-full-session-uuid",
		Slug:      "test-slug",
		Tool:      "claude-code",
		Project:   "myproject",
		StartTime: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
		Duration:  30 * time.Minute,
		TurnCount: 1,
		Messages: []parser.Message{
			{Role: "user", Text: "Hello", Timestamp: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)},
		},
	}

	got := formatJournalEntryPart(s, s.Messages, 0, 1, 1, "base", "")

	if !strings.Contains(got, `session_id: "abc12345-full-session-uuid"`) {
		t.Error("missing session_id in frontmatter")
	}
}

func TestFormatJournalEntryPart_TitleInFrontmatterAndHeading(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:        "abc12345-full-session-uuid",
		Slug:      "random-slug",
		Tool:      "claude-code",
		Project:   "myproject",
		StartTime: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
		Duration:  30 * time.Minute,
		TurnCount: 1,
		Messages: []parser.Message{
			{Role: "user", Text: "Hello", Timestamp: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)},
		},
	}

	got := formatJournalEntryPart(s, s.Messages, 0, 1, 1, "base", "Fix Authentication Bug")

	// Title should appear in frontmatter.
	if !strings.Contains(got, `title: "Fix Authentication Bug"`) {
		t.Error("missing title in frontmatter")
	}
	// Title should be used in H1 heading.
	if !strings.Contains(got, "# Fix Authentication Bug") {
		t.Error("H1 heading should use title")
	}
	// Slug should NOT be in the heading when title is provided.
	if strings.Contains(got, "# random-slug") {
		t.Error("H1 heading should use title, not slug")
	}
}

func TestFormatJournalEntryPart_NoTitleUsesSlug(t *testing.T) {
	t.Setenv("TZ", "UTC")

	s := &parser.Session{
		ID:        "abc12345-full-session-uuid",
		Slug:      "gleaming-wobbling-sutherland",
		Tool:      "claude-code",
		Project:   "myproject",
		StartTime: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
		Duration:  30 * time.Minute,
		TurnCount: 1,
		Messages: []parser.Message{
			{Role: "user", Text: "Hello", Timestamp: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)},
		},
	}

	got := formatJournalEntryPart(s, s.Messages, 0, 1, 1, "base", "")

	if !strings.Contains(got, "# gleaming-wobbling-sutherland") {
		t.Error("H1 heading should fall back to slug when no title")
	}
	// No title field in frontmatter when empty.
	if strings.Contains(got, "title:") {
		t.Error("should not have title field in frontmatter when empty")
	}
}
