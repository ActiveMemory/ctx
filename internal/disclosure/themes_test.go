//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
)

const themesTwo = "## Themes\n\n" +
	"- hooks — old hooks gist → [hooks](learnings/hooks.md)\n" +
	"- errors — error gist → [errors](learnings/errors.md)\n"

// T03: a new theme appends a bullet the parser reads back (name + link),
// on both an empty and a populated ## Themes region.
func TestWriteThemeBullet_New(t *testing.T) {
	a := disclosure.Assignment{Theme: "parsing", Slug: "parsing", Gist: "how the parser splits roots"}

	t.Run("first run creates the section", func(t *testing.T) {
		got := disclosure.WriteThemeBullet("", a, "learnings")
		want := "## Themes\n\n- parsing — how the parser splits roots → " +
			"[parsing](learnings/parsing.md)\n"
		if got != want {
			t.Fatalf("WriteThemeBullet(empty)\n got: %q\nwant: %q", got, want)
		}
		// Parses back to the same name + link.
		insp := disclosure.Inspect(
			"# Learnings\n\n"+got, disclosure.KindLearning)
		if len(insp.Themes) != 1 ||
			insp.Themes[0].Name != "parsing" ||
			insp.Themes[0].Link != "learnings/parsing.md" {
			t.Errorf("round-trip themes = %+v", insp.Themes)
		}
	})

	t.Run("append leaves existing bullets byte-identical", func(t *testing.T) {
		got := disclosure.WriteThemeBullet(themesTwo, a, "learnings")
		if !strings.Contains(got, "- hooks — old hooks gist → [hooks](learnings/hooks.md)\n") {
			t.Errorf("hooks bullet not byte-preserved:\n%q", got)
		}
		if !strings.Contains(got, "- errors — error gist → [errors](learnings/errors.md)\n") {
			t.Errorf("errors bullet not byte-preserved:\n%q", got)
		}
		if !strings.Contains(got, "- parsing — how the parser splits roots → [parsing](learnings/parsing.md)") {
			t.Errorf("new bullet missing:\n%q", got)
		}
	})
}

// T03: re-touching an existing theme replaces its bullet in place and
// leaves the other bullets byte-identical.
func TestWriteThemeBullet_Update(t *testing.T) {
	a := disclosure.Assignment{Theme: "hooks", Slug: "hooks", Gist: "NEW hooks coverage"}
	got := disclosure.WriteThemeBullet(themesTwo, a, "learnings")

	if strings.Contains(got, "old hooks gist") {
		t.Errorf("old hooks gist not replaced:\n%q", got)
	}
	if !strings.Contains(got, "- hooks — NEW hooks coverage → [hooks](learnings/hooks.md)") {
		t.Errorf("updated hooks bullet missing:\n%q", got)
	}
	// The errors bullet is untouched.
	if !strings.Contains(got, "- errors — error gist → [errors](learnings/errors.md)\n") {
		t.Errorf("errors bullet not byte-preserved:\n%q", got)
	}
	// Still exactly two bullets.
	if n := strings.Count(got, "\n- "); n != 2 {
		t.Errorf("bullet count = %d, want 2:\n%q", n, got)
	}
}
