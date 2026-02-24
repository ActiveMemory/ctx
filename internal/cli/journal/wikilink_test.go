//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"testing"
)

func TestConvertMarkdownLinks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple internal link",
			input: "[Session Title](2026-01-23-slug.md)",
			want:  "[[2026-01-23-slug|Session Title]]",
		},
		{
			name:  "internal link without extension",
			input: "[Session Title](2026-01-23-slug)",
			want:  "[[2026-01-23-slug|Session Title]]",
		},
		{
			name:  "link with parent path prefix",
			input: "[topic](../topics/caching.md)",
			want:  "[[caching|topic]]",
		},
		{
			name:  "external http link preserved",
			input: "[GitHub](https://github.com/example)",
			want:  "[GitHub](https://github.com/example)",
		},
		{
			name:  "external https link preserved",
			input: "[Docs](https://docs.example.com)",
			want:  "[Docs](https://docs.example.com)",
		},
		{
			name:  "file:// link preserved",
			input: "[View source](file:///home/user/file.md)",
			want:  "[View source](file:///home/user/file.md)",
		},
		{
			name:  "multiple links in one line",
			input: "See [A](a.md) and [B](b.md) for details",
			want:  "See [[a|A]] and [[b|B]] for details",
		},
		{
			name:  "multipart navigation link",
			input: "[← Previous](session-p1.md)",
			want:  "[[session-p1|← Previous]]",
		},
		{
			name:  "mixed internal and external",
			input: "[internal](entry.md) and [external](https://example.com)",
			want:  "[[entry|internal]] and [external](https://example.com)",
		},
		{
			name:  "no links",
			input: "Just plain text with no links",
			want:  "Just plain text with no links",
		},
		{
			name:  "deep path link",
			input: "[file](../../files/internal_cli_journal.md)",
			want:  "[[internal_cli_journal|file]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertMarkdownLinks(tt.input)
			if got != tt.want {
				t.Errorf("convertMarkdownLinks(%q)\n  got:  %q\n  want: %q",
					tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatWikilink(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		display string
		want    string
	}{
		{
			name:    "different target and display",
			target:  "2026-01-23-slug",
			display: "Session Title",
			want:    "[[2026-01-23-slug|Session Title]]",
		},
		{
			name:    "same target and display",
			target:  "caching",
			display: "caching",
			want:    "[[caching]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatWikilink(tt.target, tt.display)
			if got != tt.want {
				t.Errorf("formatWikilink(%q, %q) = %q, want %q",
					tt.target, tt.display, got, tt.want)
			}
		})
	}
}

func TestFormatWikilinkEntry(t *testing.T) {
	tests := []struct {
		name  string
		entry journalEntry
		want  string
	}{
		{
			name: "full metadata",
			entry: journalEntry{
				Filename: "2026-01-23-slug.md",
				Title:    "Session Title",
				Type:     "feature",
				Outcome:  "completed",
			},
			want: "- [[2026-01-23-slug|Session Title]] — `feature` · `completed`",
		},
		{
			name: "type only",
			entry: journalEntry{
				Filename: "2026-01-23-slug.md",
				Title:    "Session Title",
				Type:     "bugfix",
			},
			want: "- [[2026-01-23-slug|Session Title]] — `bugfix`",
		},
		{
			name: "no metadata",
			entry: journalEntry{
				Filename: "2026-01-23-slug.md",
				Title:    "Session Title",
			},
			want: "- [[2026-01-23-slug|Session Title]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatWikilinkEntry(tt.entry)
			if got != tt.want {
				t.Errorf("formatWikilinkEntry()\n  got:  %q\n  want: %q",
					got, tt.want)
			}
		})
	}
}
