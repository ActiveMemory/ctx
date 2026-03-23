//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"strings"
	"testing"
)

func TestNormalizeContent(t *testing.T) {
	tests := []struct {
		name, input    string
		fencesVerified bool
		check          func(t *testing.T, got string)
	}{
		{
			"strips tool bold",
			"🔧 **Glob: .context/journal/*.md**",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "**Glob") {
					t.Error("bold markers not stripped from tool line")
				}
				if !strings.Contains(got, "🔧 Glob:") {
					t.Error("tool prefix missing")
				}
			},
		},
		{
			"escapes glob stars",
			"pattern: src/*/main.go",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, `\*/`) {
					t.Error("glob star not escaped")
				}
			},
		},
		{
			"strips fences and escapes content",
			"```\n*.md\n```",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "```") {
					t.Error("fence markers should be stripped")
				}
			},
		},
		{
			"skips frontmatter",
			"---\ntitle: test\n---\nsome text",
			false,
			func(t *testing.T, got string) {
				if !strings.HasPrefix(got, "---\ntitle: test\n---\n") {
					t.Errorf("frontmatter mangled: %q", got)
				}
			},
		},
		{
			"does not wrap (site output is read-only)",
			"This is a very long line that exceeds eighty characters and should not be wrapped since the site output is read-only.",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "\n") {
					t.Error("NormalizeContent should not wrap lines")
				}
			},
		},
		{
			"inline code with angle brackets gets quoted",
			"the link text contains `</com` which is broken",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "`</com`") {
					t.Error("backtick code with angle bracket should be replaced")
				}
				if !strings.Contains(got, `"&lt;/com"`) {
					t.Errorf("expected quoted entity, got: %s", got)
				}
			},
		},
		{
			"tool output wrapped in pre/code",
			"### 5. Tool Output (10:30:00)\n\n# this is not a heading\n---\n<details>bad\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("tool output should be wrapped in <pre><code>")
				}
				if !strings.Contains(got, "# this is not a heading") {
					t.Error("# line should be preserved")
				}
				if !strings.Contains(got, "&lt;details&gt;bad") {
					t.Error("<details> should be HTML-escaped")
				}
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("next turn header should not be consumed")
				}
			},
		},
		{
			"boilerplate tool output stripped - empty body",
			"### 5. Tool Output (10:30:00)\n\n\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("empty tool output header should be stripped")
				}
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("next turn should survive")
				}
			},
		},
		{
			"boilerplate tool output stripped - no matches found",
			"### 5. Tool Output (10:30:00)\n\nNo matches found\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("'No matches found' tool output should be stripped")
				}
			},
		},
		{
			"non-boilerplate tool output preserved",
			"### 5. Tool Output (10:30:00)\n\nactual useful content here\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "actual useful content here") {
					t.Error("non-boilerplate content should be preserved")
				}
				if !strings.Contains(got, "<pre><code>") {
					t.Error("tool output should be wrapped in <pre><code>")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeContent(tt.input, tt.fencesVerified)
			tt.check(t, got)
		})
	}
}

func TestCollectTurnNumbersSkipsPreBlocks(t *testing.T) {
	lines := []string{
		"### 1. Assistant (10:00:00)",
		"",
		"<pre>",
		"### 800. Assistant (15:00:00)",
		"</pre>",
		"",
		"### 2. Tool Output (10:00:01)",
	}

	nums := CollectTurnNumbers(lines)

	for _, n := range nums {
		if n == 800 {
			t.Error("turn number 800 inside <pre> should be skipped")
		}
	}

	found1, found2 := false, false
	for _, n := range nums {
		if n == 1 {
			found1 = true
		}
		if n == 2 {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("expected turns 1 and 2, got %v", nums)
	}
}

func TestWrapUserTurns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, got string)
	}{
		{
			"wraps user turn in pre/code",
			"### 1. User (10:00:00)\n\nHello world\n\n### 2. Assistant (10:00:01)\n\nhi",
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("user turn should be wrapped in <pre><code>")
				}
				if !strings.Contains(got, "Hello world") {
					t.Error("user content should be preserved")
				}
				if !strings.Contains(got, "### 2. Assistant") {
					t.Error("assistant turn should survive")
				}
			},
		},
		{
			"HTML-escapes user content",
			"### 1. User (10:00:00)\n\n<script>alert('xss')</script>\n\n### 2. Assistant (10:00:01)\n\nhi",
			func(t *testing.T, got string) {
				if strings.Contains(got, "<script>") {
					t.Error("HTML should be escaped in user turn")
				}
				if !strings.Contains(got, "&lt;script&gt;") {
					t.Error("script tag should be HTML-escaped")
				}
			},
		},
		{
			"does not wrap assistant turns",
			"### 1. Assistant (10:00:00)\n\nAssistant text\n\n### 2. User (10:00:01)\n\nuser text",
			func(t *testing.T, got string) {
				if !strings.Contains(got, "Assistant text") {
					t.Error("assistant content should be preserved")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapUserTurns(tt.input)
			tt.check(t, got)
		})
	}
}
