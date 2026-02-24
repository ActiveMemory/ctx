//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"strings"
	"testing"
)

func TestContentPreview(t *testing.T) {
	tests := []struct {
		name    string
		content string
		n       int
		want    []string
	}{
		{
			name:    "simple lines",
			content: "# Heading\nFirst line\nSecond line",
			n:       3,
			want:    []string{"# Heading", "First line", "Second line"},
		},
		{
			name:    "skips frontmatter",
			content: "---\ntitle: Test\n---\n# Real Content\nBody text",
			n:       2,
			want:    []string{"# Real Content", "Body text"},
		},
		{
			name:    "skips empty lines",
			content: "Line one\n\n\nLine two\n\nLine three",
			n:       3,
			want:    []string{"Line one", "Line two", "Line three"},
		},
		{
			name:    "skips HTML comments",
			content: "<!-- comment -->\n# Heading\nContent",
			n:       2,
			want:    []string{"# Heading", "Content"},
		},
		{
			name:    "truncates long lines",
			content: strings.Repeat("a", 100),
			n:       1,
			want:    []string{strings.Repeat("a", 57) + "..."},
		},
		{
			name:    "respects limit",
			content: "one\ntwo\nthree\nfour",
			n:       2,
			want:    []string{"one", "two"},
		},
		{
			name:    "empty content",
			content: "",
			n:       3,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contentPreview(tt.content, tt.n)
			if len(got) != len(tt.want) {
				t.Fatalf("contentPreview() returned %d lines, want %d\ngot: %v", len(got), len(tt.want), got)
			}
			for i, line := range got {
				if line != tt.want[i] {
					t.Errorf("line[%d] = %q, want %q", i, line, tt.want[i])
				}
			}
		})
	}
}
