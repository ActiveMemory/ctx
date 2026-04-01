//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import "testing"

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
			got := FenceForContent(tt.content)
			if got != tt.want {
				t.Errorf("FenceForContent(%q) = %q, want %q", tt.content, got, tt.want)
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
			got := StripLineNumbers(tt.content)
			if got != tt.want {
				t.Errorf("StripLineNumbers() = %q, want %q", got, tt.want)
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
			name: "single reminder",
			content: "before " +
				"<system-reminder>reminder text</system-reminder> after",
			wantClean:     "before  after",
			wantReminders: 1,
		},
		{
			name: "multiple reminders",
			content: "<system-reminder>first</system-reminder>" +
				" middle " +
				"<system-reminder>second</system-reminder>",
			wantClean:     " middle ",
			wantReminders: 2,
		},
		{
			name: "multiline reminder",
			content: "text\n<system-reminder>\n" +
				"multiline\nreminder\n" +
				"</system-reminder>\nmore",
			wantClean:     "text\n\nmore",
			wantReminders: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClean, gotReminders := ExtractSystemReminders(tt.content)
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
			got := NormalizeCodeFences(tt.content)
			if got != tt.want {
				t.Errorf("NormalizeCodeFences() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}
