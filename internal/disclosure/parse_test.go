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

	// Convention fixtures (M4): the same preamble | staging | ## Themes
	// layout as the entry kinds. Staging holds curated "## " prose
	// sections; the retired model nested "### " under a "## Recent".
	conventionMigrated = "# Conventions\n\n<!-- guide -->\n\n" +
		"## Error handling\n\nwrap with %w.\n\n" +
		"## Naming\n\nname files for their concern.\n\n" +
		"## Themes\n\n- naming — file naming → [naming](conventions/naming.md)\n"

	conventionUnmigrated = "# Conventions\n\n<!-- guide -->\n\n" +
		"## a convention\n\nprose.\n"

	// A "## " line inside the preamble's HTML comment is an illustrative
	// example in the authoring guide, not a staged section.
	conventionCommentedExample = "# Conventions\n\n" +
		"<!--\nAdd sections like:\n\n## Section Title\n\nprose.\n-->\n\n" +
		"## Real Section\n\nreal prose.\n"
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
		{"convention commented example", conventionCommentedExample, disclosure.KindConvention},
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

// M4/T07: a convention root splits on the same preamble | staging |
// ## Themes layout as the entry kinds — staging is the "## " sections
// above the themes region, not a "## Recent" block below it.
func TestParse_ConventionKind(t *testing.T) {
	r := disclosure.Parse(conventionMigrated, disclosure.KindConvention)
	if !r.HasThemes {
		t.Fatal("HasThemes = false, want true")
	}
	if !strings.HasPrefix(r.Staging, "## Error handling") {
		t.Errorf("staging must start at the first section; staging=%q", r.Staging)
	}
	if !strings.Contains(r.Staging, "## Naming") {
		t.Errorf("second section not in staging; staging=%q", r.Staging)
	}
	if strings.Contains(r.Staging, "## Themes") {
		t.Errorf("staging leaked the themes section; staging=%q", r.Staging)
	}
	if !strings.HasPrefix(r.Preamble, "# Conventions") ||
		strings.Contains(r.Preamble, "## Error handling") {
		t.Errorf("preamble wrong; preamble=%q", r.Preamble)
	}
	if len(r.Themes) != 1 || r.Themes[0].Link != "conventions/naming.md" {
		t.Errorf("themes parse wrong: %+v", r.Themes)
	}
}

// M4/T05/T09 (C2): the convention enumerator reports one entry per
// curated section, with title identity and no timestamp.
func TestStagedEntries_Convention(t *testing.T) {
	r := disclosure.Parse(conventionMigrated, disclosure.KindConvention)
	got := disclosure.StagedEntries(r)
	if len(got) != 2 {
		t.Fatalf("StagedEntries = %d, want 2: %+v", len(got), got)
	}
	want := []string{"Error handling", "Naming"}
	for i, w := range want {
		if got[i].Title != w {
			t.Errorf("entry %d title = %q, want %q", i, got[i].Title, w)
		}
		if got[i].Timestamp != "" {
			t.Errorf("entry %d timestamp = %q, want empty (title identity)",
				i, got[i].Timestamp)
		}
	}
}

// M4/T09 (C7): a "## " line inside an HTML comment is an authoring
// example, not a staged section — the comment-skip carried over from M3
// must hold for the convention prefix, which collides far more readily
// than "## [".
func TestParse_ConventionCommentSkip(t *testing.T) {
	r := disclosure.Parse(conventionCommentedExample, disclosure.KindConvention)
	if !strings.HasPrefix(r.Staging, "## Real Section") {
		t.Errorf("staging must start at the real section; staging=%q", r.Staging)
	}
	if strings.Contains(r.Staging, "Add sections like") {
		t.Errorf("comment leaked into staging; staging=%q", r.Staging)
	}
	entries := disclosure.StagedEntries(r)
	if len(entries) != 1 || entries[0].Title != "Real Section" {
		t.Errorf("StagedEntries = %+v, want just Real Section", entries)
	}
}

// A "## [" line inside an HTML comment (a knowledge file's format guide —
// e.g. DECISIONS.md's "## [YYYY-MM-DD] Decision Title" example) must not be
// mistaken for a staging entry. Regression: after a fold empties the real
// staging, such a commented example would otherwise leave a non-empty,
// unparsable staging zone and trip Validate (ErrStagingUnparsable).
func TestParse_CommentedEntryExampleNotStaging(t *testing.T) {
	const root = "# Decisions\n\n<!-- DECISION FORMATS\n\n" +
		"## [YYYY-MM-DD] Decision Title\n\n**Status**: Accepted\n\n-->\n\n" +
		"## Themes\n\n- sec — security → [sec](decisions/sec.md)\n"

	r := disclosure.Parse(root, disclosure.KindDecision)
	if strings.TrimSpace(r.Staging) != "" {
		t.Errorf("Staging = %q, want empty (commented example belongs in preamble)",
			r.Staging)
	}
	if !strings.Contains(r.Preamble, "## [YYYY-MM-DD]") {
		t.Errorf("commented example missing from preamble; preamble=%q", r.Preamble)
	}
	if r.Reconstruct() != root {
		t.Error("round-trip mismatch after comment-aware parse")
	}
	if err := disclosure.Validate(r); err != nil {
		t.Errorf("Validate = %v, want nil (commented example is not real staging)", err)
	}
}

// A "## Themes" line inside an HTML comment is an example, not a second
// themes section: Validate must not read it as a duplicate (ErrMultipleThemes).
func TestParse_CommentedThemesNotCounted(t *testing.T) {
	const root = "# Learnings\n\n<!-- example: ## Themes goes here -->\n\n" +
		"## Themes\n\n- hooks — hook mechanics → [hooks](learnings/hooks.md)\n"

	r := disclosure.Parse(root, disclosure.KindLearning)
	if !r.HasThemes {
		t.Fatal("HasThemes = false, want true (one real ## Themes)")
	}
	if err := disclosure.Validate(r); err != nil {
		t.Errorf("Validate = %v, want nil (commented ## Themes is not a duplicate)", err)
	}
}
