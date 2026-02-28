//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
)

func TestStripMkDocs_Frontmatter(t *testing.T) {
	input := "---\ntitle: Test\nicon: lucide/flame\n---\n\n# Hello\n"
	got := StripMkDocs(input)
	if strings.Contains(got, "title: Test") {
		t.Error("frontmatter was not stripped")
	}
	if !strings.Contains(got, "# Hello") {
		t.Error("content after frontmatter was lost")
	}
}

func TestStripMkDocs_Images(t *testing.T) {
	input := "# Title\n\n![banner](images/banner.png)\n\nSome text.\n"
	got := StripMkDocs(input)
	if strings.Contains(got, "![banner]") {
		t.Error("image line was not stripped")
	}
	if !strings.Contains(got, "Some text.") {
		t.Error("text after image was lost")
	}
}

func TestStripMkDocs_Admonitions(t *testing.T) {
	input := `!!! note "Important Note"
    This is the body.
    Second line.

Normal text.
`
	got := StripMkDocs(input)
	if !strings.Contains(got, `> **Important Note**`) {
		t.Errorf("admonition title not converted, got:\n%s", got)
	}
	if !strings.Contains(got, "> This is the body.") {
		t.Errorf("admonition body not dedented to blockquote, got:\n%s", got)
	}
	if !strings.Contains(got, "> Second line.") {
		t.Errorf("admonition second line not dedented, got:\n%s", got)
	}
	if !strings.Contains(got, "Normal text.") {
		t.Error("text after admonition was lost")
	}
}

func TestStripMkDocs_AdmonitionNoTitle(t *testing.T) {
	input := "!!! warning\n    Body here.\n"
	got := StripMkDocs(input)
	if strings.Contains(got, "!!!") {
		t.Error("admonition marker was not stripped")
	}
	if !strings.Contains(got, "> Body here.") {
		t.Errorf("admonition body not converted, got:\n%s", got)
	}
}

func TestStripMkDocs_Tabs(t *testing.T) {
	input := `=== "Without ctx"

    Some content here.
    More content.

=== "With ctx"

    Better content.

Normal text.
`
	got := StripMkDocs(input)
	if !strings.Contains(got, "**Without ctx**") {
		t.Errorf("tab title not converted, got:\n%s", got)
	}
	if !strings.Contains(got, "Some content here.") {
		t.Error("tab body not dedented")
	}
	if strings.Contains(got, "    Some content here.") {
		t.Error("tab body still has 4-space indent")
	}
	if !strings.Contains(got, "**With ctx**") {
		t.Error("second tab title not converted")
	}
	if !strings.Contains(got, "Normal text.") {
		t.Error("text after tabs was lost")
	}
}

func TestStripMkDocs_RelativeLinks(t *testing.T) {
	input := "[Getting Started](getting-started.md) and [About](../home/about.md)\n"
	got := StripMkDocs(input)
	if strings.Contains(got, ".md") {
		t.Errorf("relative .md link not stripped, got:\n%s", got)
	}
	if !strings.Contains(got, "Getting Started") {
		t.Error("link text was lost")
	}
	if !strings.Contains(got, "About") {
		t.Error("second link text was lost")
	}
}

func TestStripMkDocs_PreservesExternalLinks(t *testing.T) {
	input := "[ctx site](https://ctx.ist) stays.\n"
	got := StripMkDocs(input)
	if !strings.Contains(got, "[ctx site](https://ctx.ist)") {
		t.Errorf("external link was modified, got:\n%s", got)
	}
}

func TestStripMkDocs_PreservesCodeBlocks(t *testing.T) {
	input := "```text\n{} --> what\nctx --> why\n```\n"
	got := StripMkDocs(input)
	if !strings.Contains(got, "ctx --> why") {
		t.Errorf("code block content was lost, got:\n%s", got)
	}
}

func TestStripMkDocs_EmbeddedManifesto(t *testing.T) {
	content, loadErr := assets.WhyDoc("manifesto")
	if loadErr != nil {
		t.Fatalf("failed to load embedded manifesto: %v", loadErr)
	}

	got := StripMkDocs(string(content))

	// Should not contain MkDocs artifacts.
	if strings.Contains(got, "---\ntitle:") {
		t.Error("frontmatter not stripped from manifesto")
	}
	if strings.Contains(got, "![ctx]") {
		t.Error("image not stripped from manifesto")
	}
	if strings.Contains(got, "!!! ") {
		t.Error("admonition markers not stripped from manifesto")
	}

	// Should contain key content.
	if !strings.Contains(got, "Manifesto") {
		t.Error("manifesto title missing")
	}
	if !strings.Contains(got, "Creation, not code") {
		t.Error("manifesto opening missing")
	}
}
