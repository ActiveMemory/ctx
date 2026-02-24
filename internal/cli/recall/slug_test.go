//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

func TestSlugifyTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"simple title", "Fix Authentication Bug", "fix-authentication-bug"},
		{"with punctuation", "Hello, World! How's it going?", "hello-world-how-s-it-going"},
		{"leading/trailing spaces", "  hello world  ", "hello-world"},
		{"unicode characters", "café résumé naïve", "caf-r-sum-na-ve"},
		{"unicode CJK", "修复认证错误", ""},
		{"single word", "refactor", "refactor"},
		{"numbers", "Add OAuth2 support for v3 API", "add-oauth2-support-for-v3-api"},
		{"consecutive special chars", "hello---world!!!test", "hello-world-test"},
		{"all punctuation", "!@#$%^&*()", ""},
		{"truncated FirstUserMsg suffix", "implement the feature...", "implement-the-feature"},
		{
			"very long title truncates on word boundary",
			"This is an extremely long title that should be truncated on a word boundary somewhere around fifty characters",
			"this-is-an-extremely-long-title-that-should-be",
		},
		{
			"long title no good break point",
			strings.Repeat("a", 60),
			strings.Repeat("a", 50),
		},
		{"just hey", "hey", "hey"},
		{"mixed case", "Fix README.md Formatting Issues", "fix-readme-md-formatting-issues"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugifyTitle(tt.input)
			if got != tt.want {
				t.Errorf("slugifyTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
			// Slug must not exceed max length.
			if len(got) > slugMaxLen {
				t.Errorf("slugifyTitle(%q) length %d exceeds max %d", tt.input, len(got), slugMaxLen)
			}
		})
	}
}

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "Fix Authentication Bug", "Fix Authentication Bug"},
		{"with newlines", "Implement plan:\n\n# Title\n\nDetails", "Implement plan: # Title Details"},
		{"with tabs", "Hello\tworld", "Hello world"},
		{"leading/trailing whitespace", "  hello  ", "hello"},
		{"truncation suffix", "some long text...", "some long text"},
		{"consecutive spaces", "hello    world", "hello world"},
		{"empty", "", ""},
		{
			"long title truncated on word boundary",
			"We are debugging the new journal enrichment and creation flow and the rendering still breaks around line 37",
			"We are debugging the new journal enrichment and creation flow and the",
		},
		{
			"title at exactly 75 chars is not truncated",
			strings.Repeat("x ", 37) + "y",
			strings.Repeat("x ", 37) + "y",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanTitle(tt.input)
			if got != tt.want {
				t.Errorf("cleanTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTitleSlug_FallbackHierarchy(t *testing.T) {
	tests := []struct {
		name          string
		session       *parser.Session
		existingTitle string
		wantSlug      string
		wantTitle     string
	}{
		{
			name: "prefers existingTitle",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "random-slug",
				FirstUserMsg: "implement auth",
			},
			existingTitle: "Fix Authentication Bug",
			wantSlug:      "fix-authentication-bug",
			wantTitle:     "Fix Authentication Bug",
		},
		{
			name: "falls back to FirstUserMsg",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "random-slug",
				FirstUserMsg: "implement the new auth system",
			},
			existingTitle: "",
			wantSlug:      "implement-the-new-auth-system",
			wantTitle:     "implement the new auth system",
		},
		{
			name: "FirstUserMsg with newlines gets cleaned",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "random-slug",
				FirstUserMsg: "Implement plan:\n\n# Title...",
			},
			existingTitle: "",
			wantSlug:      "implement-plan-title",
			wantTitle:     "Implement plan: # Title",
		},
		{
			name: "falls back to Slug",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "gleaming-wobbling-sutherland",
				FirstUserMsg: "",
			},
			existingTitle: "",
			wantSlug:      "gleaming-wobbling-sutherland",
			wantTitle:     "",
		},
		{
			name: "falls back to short ID",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "",
				FirstUserMsg: "",
			},
			existingTitle: "",
			wantSlug:      "abc12345",
			wantTitle:     "",
		},
		{
			name: "existingTitle all punctuation falls through",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "fallback-slug",
				FirstUserMsg: "hello world",
			},
			existingTitle: "!@#$%",
			wantSlug:      "hello-world",
			wantTitle:     "hello world",
		},
		{
			name: "FirstUserMsg all punctuation falls to Slug",
			session: &parser.Session{
				ID:           "abc12345-full-uuid",
				Slug:         "fallback-slug",
				FirstUserMsg: "...",
			},
			existingTitle: "",
			wantSlug:      "fallback-slug",
			wantTitle:     "",
		},
		{
			name: "short ID when ID is short",
			session: &parser.Session{
				ID:           "abc",
				Slug:         "",
				FirstUserMsg: "",
			},
			existingTitle: "",
			wantSlug:      "abc",
			wantTitle:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlug, gotTitle := titleSlug(tt.session, tt.existingTitle)
			if gotSlug != tt.wantSlug {
				t.Errorf("titleSlug() slug = %q, want %q", gotSlug, tt.wantSlug)
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("titleSlug() title = %q, want %q", gotTitle, tt.wantTitle)
			}
		})
	}
}
