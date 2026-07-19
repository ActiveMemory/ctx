//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
)

// T01: KindFor maps the three canonical filenames and rejects others.
func TestKindFor(t *testing.T) {
	cases := []struct {
		name    string
		wantOK  bool
		wantStr string
	}{
		{"LEARNINGS.md", true, "learning"},
		{"DECISIONS.md", true, "decision"},
		{"CONVENTIONS.md", true, "convention"},
		{"README.md", false, ""},
		{"learnings.md", false, ""}, // exact match only
		{"", false, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k, ok := disclosure.KindFor(tc.name)
			if ok != tc.wantOK {
				t.Fatalf("KindFor(%q) ok = %v, want %v", tc.name, ok, tc.wantOK)
			}
			if ok && k.String() != tc.wantStr {
				t.Errorf("KindFor(%q).String() = %q, want %q",
					tc.name, k.String(), tc.wantStr)
			}
		})
	}
}

// T01 (pd-m3): ThemeDir maps the digestible entry kinds to their theme
// subdirectory and refuses convention/unknown so the mover never writes
// to a guessed path.
func TestThemeDir(t *testing.T) {
	cases := []struct {
		name   string
		kind   disclosure.Kind
		wantOK bool
		want   string
	}{
		{"learning", disclosure.KindLearning, true, "learnings"},
		{"decision", disclosure.KindDecision, true, "decisions"},
		{"convention", disclosure.KindConvention, false, ""},
		{"unknown", disclosure.Kind(99), false, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := disclosure.ThemeDir(tc.kind)
			if ok != tc.wantOK {
				t.Fatalf("ThemeDir(%v) ok = %v, want %v", tc.kind, ok, tc.wantOK)
			}
			if got != tc.want {
				t.Errorf("ThemeDir(%v) = %q, want %q", tc.kind, got, tc.want)
			}
		})
	}
}

// T02: StagedEntries lists the staging zone's entries in order, and is
// empty when staging holds none.
func TestStagedEntries(t *testing.T) {
	populated := disclosure.Parse(entryMigratedPopulated, disclosure.KindLearning)
	got := disclosure.StagedEntries(populated)
	if len(got) != 1 {
		t.Fatalf("StagedEntries = %d entries, want 1", len(got))
	}
	if got[0].Timestamp != "2026-07-15-120000" ||
		got[0].Title != "a staged entry" {
		t.Errorf("StagedEntries[0] = %+v, want ts/title of the staged entry", got[0])
	}

	empty := disclosure.Parse(entryMigratedEmpty, disclosure.KindLearning)
	if s := disclosure.StagedEntries(empty); s != nil {
		t.Errorf("StagedEntries(empty staging) = %+v, want nil", s)
	}
}

// T03: Inspect assembles kind + staging + themes consistent with Parse.
func TestInspect(t *testing.T) {
	insp := disclosure.Inspect(entryMigratedPopulated, disclosure.KindLearning)

	if insp.Kind != "learning" {
		t.Errorf("Inspect.Kind = %q, want learning", insp.Kind)
	}
	if len(insp.Staging) != 1 || insp.Staging[0].Title != "a staged entry" {
		t.Errorf("Inspect.Staging = %+v, want the one staged entry", insp.Staging)
	}
	if len(insp.Themes) != 1 || insp.Themes[0].Link != "learnings/hooks.md" {
		t.Errorf("Inspect.Themes = %+v, want the one theme", insp.Themes)
	}
}
