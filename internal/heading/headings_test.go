//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heading

import (
	"slices"
	"testing"

	"github.com/ActiveMemory/ctx/internal/entity"
)

func texts(hs []entity.Heading) []string {
	out := make([]string, len(hs))
	for i, h := range hs {
		out[i] = h.Text
	}
	return out
}

func eq(a, b []string) bool {
	return slices.Equal(a, b)
}

func TestHeadings_DefaultLevelTwoOnly(t *testing.T) {
	content := "# Title\n\n## Alpha\n\ntext\n\n### Sub of Alpha\n\n## Beta\n"
	got := texts(Headings(content, 2))
	want := []string{"Alpha", "Beta"}
	if !eq(got, want) {
		t.Fatalf("depth 2: got %v want %v", got, want)
	}
}

func TestHeadings_DepthThreeIncludesLevelThree(t *testing.T) {
	content := "# Title\n\n## Alpha\n\n### Sub of Alpha\n\n## Beta\n"
	got := texts(Headings(content, 3))
	want := []string{"Alpha", "Sub of Alpha", "Beta"}
	if !eq(got, want) {
		t.Fatalf("depth 3: got %v want %v", got, want)
	}
}

func TestHeadings_UntimestampedPhasesMatch(t *testing.T) {
	// TASKS.md-style headings carry no timestamp; the generic matcher must
	// still project them (unlike ParseHeaders).
	content := "# Tasks\n\n## Phase 1: Foundation\n\n- [ ] do a thing\n\n## Phase 2: Ship\n"
	got := texts(Headings(content, 2))
	want := []string{"Phase 1: Foundation", "Phase 2: Ship"}
	if !eq(got, want) {
		t.Fatalf("phases: got %v want %v", got, want)
	}
}

func TestHeadings_FencedHeadingIgnored(t *testing.T) {
	content := "## Real\n\n```\n## not a heading\n```\n\n## Also Real\n"
	got := texts(Headings(content, 3))
	want := []string{"Real", "Also Real"}
	if !eq(got, want) {
		t.Fatalf("fenced: got %v want %v", got, want)
	}
}

func TestHeadings_HTMLCommentHeadingsIgnored(t *testing.T) {
	// The `<!-- DECISION FORMATS ... -->` legend embeds example ## headings
	// that are documentation, not entries; they must not be projected.
	content := "# Decisions\n\n<!-- FORMATS\n\n## Quick Format\n\n## Full Format\n-->\n\n## [2026-07-06-120000] Real Decision\n"
	got := texts(Headings(content, 3))
	want := []string{"[2026-07-06-120000] Real Decision"}
	if !eq(got, want) {
		t.Fatalf("html comment: got %v want %v", got, want)
	}
}

func TestHeadings_InlineCommentNotABlock(t *testing.T) {
	// A self-contained inline comment must not swallow following headings.
	content := "<!-- INDEX:START --><!-- INDEX:END -->\n\n## Real\n"
	got := texts(Headings(content, 2))
	want := []string{"Real"}
	if !eq(got, want) {
		t.Fatalf("inline comment: got %v want %v", got, want)
	}
}

func TestHeadings_EmptyAndNoHeading(t *testing.T) {
	if got := Headings("", 2); len(got) != 0 {
		t.Fatalf("empty: got %v want none", got)
	}
	if got := Headings("just prose\nno headings here\n", 2); len(got) != 0 {
		t.Fatalf("no-heading: got %v want none", got)
	}
}

func TestHeadings_StaleIndexBlockIgnored(t *testing.T) {
	// A residual INDEX comment block is not a heading and must not appear.
	content := "# Decisions\n\n<!-- INDEX:START -->\n| Date | Decision |\n|----|--------|\n| 2026-07-06 | Something |\n<!-- INDEX:END -->\n\n## [2026-07-06-120000] Something\n"
	got := texts(Headings(content, 2))
	want := []string{"[2026-07-06-120000] Something"}
	if !eq(got, want) {
		t.Fatalf("stale block: got %v want %v", got, want)
	}
}

func TestHeadings_DepthBelowTwoClampsToTwo(t *testing.T) {
	content := "# Title\n\n## Alpha\n"
	got := texts(Headings(content, 0))
	want := []string{"Alpha"}
	if !eq(got, want) {
		t.Fatalf("clamp: got %v want %v", got, want)
	}
}

func TestHeadings_LevelRecorded(t *testing.T) {
	content := "## Alpha\n### Beta\n"
	got := Headings(content, 3)
	if len(got) != 2 || got[0].Level != 2 || got[1].Level != 3 {
		t.Fatalf("levels: got %+v", got)
	}
}
