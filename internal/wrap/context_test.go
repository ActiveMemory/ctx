//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package wrap

import (
	"strings"
	"testing"
)

func TestContextFile_TaskLine(t *testing.T) {
	input := "- [ ] Rename 8 skills to ctx-domain-action naming convention across all skill files and update cross-references #priority:high #session:a92cadca #branch:main #commit:c4c53c7a #added:2026-04-06-212611"
	got := ContextFile(input, 80)
	lines := strings.Split(got, "\n")

	if len(lines) < 2 {
		t.Fatalf("expected wrapping, got %d lines", len(lines))
	}
	if !strings.HasPrefix(lines[0], "- [ ] ") {
		t.Error("first line should keep checkbox prefix")
	}
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "  ") {
			t.Errorf(
				"continuation line should have 2-space indent: %q",
				line,
			)
		}
	}
}

func TestContextFile_SkipsHeadings(t *testing.T) {
	input := "## [2026-04-06-123456] This is a very long heading that exceeds eighty characters and should not be wrapped at all by the formatter"
	got := ContextFile(input, 80)

	if got != input {
		t.Errorf("heading should not be wrapped:\n  got:  %q\n  want: %q", got, input)
	}
}

func TestContextFile_SkipsTables(t *testing.T) {
	input := "| Column A | Column B | Column C | Column D | Column E | Column F | Column G | Column H |"
	got := ContextFile(input, 80)

	if got != input {
		t.Errorf("table row should not be wrapped")
	}
}

func TestContextFile_SkipsFrontmatter(t *testing.T) {
	input := "---\ntitle: This is a very long frontmatter value that exceeds eighty characters and should not be wrapped\n---\nShort line"
	got := ContextFile(input, 80)

	if !strings.Contains(got, "title: This is a very long frontmatter") {
		t.Error("frontmatter should not be wrapped")
	}
}

func TestContextFile_SkipsHTMLComments(t *testing.T) {
	input := "<!-- This is a very long HTML comment that exceeds eighty characters and should not be wrapped by the formatter at all -->"
	got := ContextFile(input, 80)

	if got != input {
		t.Error("HTML comment should not be wrapped")
	}
}

func TestContextFile_Idempotent(t *testing.T) {
	input := "- [ ] Rename 8 skills to ctx-domain-action naming convention across all skill files and update cross-references #priority:high #added:2026-04-06-212611\n\n## Heading\n\n| Col A | Col B |\n"
	first := ContextFile(input, 80)
	second := ContextFile(first, 80)

	if first != second {
		t.Errorf(
			"not idempotent:\n  first:  %q\n  second: %q",
			first, second,
		)
	}
}

func TestContextFile_ShortLinesUnchanged(t *testing.T) {
	input := "- [ ] Short task #added:2026-04-06\n- [x] Done task"
	got := ContextFile(input, 80)

	if got != input {
		t.Errorf("short lines should not change:\n  got:  %q\n  want: %q", got, input)
	}
}

func TestContextFile_ConventionLine(t *testing.T) {
	input := "- Use camelCase for import aliases in all Go files across the project, matching the existing convention established in the codebase"
	got := ContextFile(input, 80)
	lines := strings.Split(got, "\n")

	if len(lines) < 2 {
		t.Fatal("expected wrapping for long convention line")
	}
	if !strings.HasPrefix(lines[1], "  ") {
		t.Errorf(
			"continuation should have 2-space indent: %q",
			lines[1],
		)
	}
}

func TestContextFile_ParagraphText(t *testing.T) {
	input := "**Context**: We decided to use the new authentication middleware because the old one was flagged by legal for storing session tokens in a way that does not meet compliance requirements."
	got := ContextFile(input, 80)
	lines := strings.Split(got, "\n")

	if len(lines) < 2 {
		t.Fatal("expected wrapping for long paragraph")
	}
}

func TestContextFile_MultilineHTMLComment(t *testing.T) {
	input := "<!-- START\nThis is a very long line inside a multiline HTML comment block that exceeds eighty characters and should be preserved\nEND -->\nRegular line"
	got := ContextFile(input, 80)

	if !strings.Contains(got, "This is a very long line inside a multiline") {
		t.Error("content inside HTML comment should not be wrapped")
	}
}

func TestListItem_NoBreakMidWord(t *testing.T) {
	input := "- [ ] https://very-long-url-that-has-no-spaces-and-exceeds-eighty-characters-by-far-and-cannot-be-broken.example.com/path"
	got := ListItem(input, 80)

	// URL may end up on a continuation line, but should never be
	// split mid-word.
	for _, line := range got {
		if strings.Contains(line, "long-url") &&
			!strings.Contains(line, ".example.com/path") {
			t.Errorf("URL was split mid-word: %q", line)
		}
	}
}
