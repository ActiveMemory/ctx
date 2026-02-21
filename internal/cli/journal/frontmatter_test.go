//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"strings"
	"testing"
)

func TestTransformFrontmatter(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		sourcePath string
		wantTags   bool // expect "tags:" instead of "topics:"
		wantAlias  bool // expect "aliases:" with title
		wantSource bool // expect "source_file:" field
		wantUnmod  bool // expect content unchanged
	}{
		{
			name:       "topics renamed to tags",
			content:    "---\ntitle: Test Session\ndate: 2026-01-23\ntopics:\n  - caching\n  - auth\n---\nBody content\n",
			sourcePath: ".context/journal/2026-01-23-test.md",
			wantTags:   true,
			wantAlias:  true,
			wantSource: true,
		},
		{
			name:      "no frontmatter passthrough",
			content:   "# Just a heading\n\nBody content\n",
			wantUnmod: true,
		},
		{
			name:      "incomplete frontmatter passthrough",
			content:   "---\ntitle: Incomplete\n",
			wantUnmod: true,
		},
		{
			name:       "preserves type and outcome",
			content:    "---\ntitle: Feature Work\ndate: 2026-02-01\ntype: feature\noutcome: completed\ntopics:\n  - auth\n---\nBody\n",
			sourcePath: ".context/journal/entry.md",
			wantTags:   true,
			wantAlias:  true,
			wantSource: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformFrontmatter(tt.content, tt.sourcePath)

			if tt.wantUnmod {
				if got != tt.content {
					t.Errorf("expected unmodified content, got:\n%s", got)
				}
				return
			}

			if tt.wantTags {
				if strings.Contains(got, "topics:") {
					t.Error("output still contains 'topics:' — should be 'tags:'")
				}
				if !strings.Contains(got, "tags:") {
					t.Error("output missing 'tags:' field")
				}
			}

			if tt.wantAlias {
				if !strings.Contains(got, "aliases:") {
					t.Error("output missing 'aliases:' field")
				}
			}

			if tt.wantSource {
				if !strings.Contains(got, "source_file:") {
					t.Error("output missing 'source_file:' field")
				}
				if !strings.Contains(got, tt.sourcePath) {
					t.Errorf("output missing source path %q", tt.sourcePath)
				}
			}

			// Body should be preserved after frontmatter
			if !strings.Contains(got, "Body") {
				t.Error("body content was lost during transformation")
			}

			// Should still have frontmatter delimiters
			if !strings.HasPrefix(got, "---\n") {
				t.Error("output missing opening frontmatter delimiter")
			}
		})
	}
}

func TestTransformFrontmatterPreservesBody(t *testing.T) {
	content := "---\ntitle: Test\ndate: 2026-01-23\ntopics:\n  - go\n---\n# Heading\n\nParagraph one.\n\nParagraph two.\n"
	got := transformFrontmatter(content, "source.md")

	// Find the body after the closing ---
	parts := strings.SplitN(got, "---\n", 3)
	if len(parts) < 3 {
		t.Fatalf("expected 3 parts (open, fm, body), got %d", len(parts))
	}

	body := parts[2]
	if !strings.Contains(body, "# Heading") {
		t.Error("body missing heading")
	}
	if !strings.Contains(body, "Paragraph one.") {
		t.Error("body missing paragraph one")
	}
	if !strings.Contains(body, "Paragraph two.") {
		t.Error("body missing paragraph two")
	}
}

func TestExtractStringSlice(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		key  string
		want int // expected length, -1 for nil
	}{
		{
			name: "string slice",
			m:    map[string]any{"tags": []any{"a", "b", "c"}},
			key:  "tags",
			want: 3,
		},
		{
			name: "missing key",
			m:    map[string]any{"other": "value"},
			key:  "tags",
			want: -1,
		},
		{
			name: "non-slice value",
			m:    map[string]any{"tags": "single"},
			key:  "tags",
			want: -1,
		},
		{
			name: "empty slice",
			m:    map[string]any{"tags": []any{}},
			key:  "tags",
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractStringSlice(tt.m, tt.key)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			} else if len(got) != tt.want {
				t.Errorf("expected len %d, got %d", tt.want, len(got))
			}
		})
	}
}
