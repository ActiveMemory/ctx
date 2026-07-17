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

// Entry-kind fixtures (LEARNINGS/DECISIONS): staging above ## Themes.
const (
	entryMigratedPopulated = "# Learnings\n\n<!-- guide -->\n\n" +
		"## [2026-07-15-120000] a staged entry\n\n**Context**: x.\n\n---\n\n" +
		"## Themes\n\n- hooks — hook mechanics → [hooks](learnings/hooks.md)\n"

	entryMigratedEmpty = "# Learnings\n\n<!-- guide -->\n\n" +
		"## Themes\n\n- hooks — hook mechanics → [hooks](learnings/hooks.md)\n"

	entryUnmigrated = "# Learnings\n\n<!-- guide -->\n\n" +
		"## [2026-07-15-120000] a staged entry\n\n**Context**: x.\n"

	conventionMigrated = "# Conventions\n\n<!-- guide -->\n\n" +
		"## Themes\n\n- naming — file naming → [naming](conventions/naming.md)\n\n" +
		"## Recent\n\n### a recent convention\n\nprose.\n"

	conventionUnmigrated = "# Conventions\n\n<!-- guide -->\n\n" +
		"### a convention\n\nprose.\n"
)

// T04/T05: Parse must round-trip every shape byte-for-byte — nothing is
// normalized, so the mover (M2) gets exact bytes.
func TestParse_RoundTrip(t *testing.T) {
	cases := []struct {
		name    string
		content string
		kind    disclosure.Kind
	}{
		{"entry migrated populated", entryMigratedPopulated, disclosure.KindLearning},
		{"entry migrated empty staging", entryMigratedEmpty, disclosure.KindLearning},
		{"entry un-migrated", entryUnmigrated, disclosure.KindDecision},
		{"convention migrated", conventionMigrated, disclosure.KindConvention},
		{"convention un-migrated", conventionUnmigrated, disclosure.KindConvention},
		{"empty", "", disclosure.KindLearning},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := disclosure.Parse(tc.content, tc.kind).Reconstruct()
			if got != tc.content {
				t.Errorf("round-trip mismatch\n got: %q\nwant: %q", got, tc.content)
			}
		})
	}
}

// T04: entry-kind split places entries in staging and gists in themes.
func TestParse_EntryKind(t *testing.T) {
	r := disclosure.Parse(entryMigratedPopulated, disclosure.KindLearning)
	if !r.HasThemes {
		t.Fatal("HasThemes = false, want true (## Themes present)")
	}
	if !strings.Contains(r.Staging, "a staged entry") {
		t.Errorf("staging missing the entry; staging=%q", r.Staging)
	}
	if strings.Contains(r.Staging, "## Themes") {
		t.Errorf("staging leaked the themes section; staging=%q", r.Staging)
	}
	if len(r.Themes) != 1 || r.Themes[0].Link != "learnings/hooks.md" {
		t.Errorf("themes parse wrong: %+v", r.Themes)
	}
}

// T04: an un-migrated entry root has empty themes and all entries staged.
func TestParse_EntryKindUnmigrated(t *testing.T) {
	r := disclosure.Parse(entryUnmigrated, disclosure.KindLearning)
	if r.HasThemes {
		t.Error("HasThemes = true, want false (no ## Themes)")
	}
	if !strings.Contains(r.Staging, "a staged entry") {
		t.Errorf("entry not in staging; staging=%q", r.Staging)
	}
	if r.ThemesRaw != "" {
		t.Errorf("ThemesRaw = %q, want empty", r.ThemesRaw)
	}
}

// T05: conventions split has themes between ## Themes and ## Recent, and
// staging is the ## Recent region.
func TestParse_ConventionKind(t *testing.T) {
	r := disclosure.Parse(conventionMigrated, disclosure.KindConvention)
	if !r.HasThemes {
		t.Fatal("HasThemes = false, want true")
	}
	if !strings.HasPrefix(strings.TrimSpace(r.Staging), "## Recent") {
		t.Errorf("staging must start at ## Recent; staging=%q", r.Staging)
	}
	if !strings.Contains(r.Staging, "a recent convention") {
		t.Errorf("recent section not in staging; staging=%q", r.Staging)
	}
	if strings.Contains(r.ThemesRaw, "## Recent") {
		t.Errorf("themes region leaked ## Recent; themesRaw=%q", r.ThemesRaw)
	}
	if len(r.Themes) != 1 || r.Themes[0].Link != "conventions/naming.md" {
		t.Errorf("themes parse wrong: %+v", r.Themes)
	}
}
