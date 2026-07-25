//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert_test

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/insert"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
)

// Layout proof for specs/progressive-disclosure.md (plan pd-m4, T16/T17):
// CONVENTIONS. A convention is a bullet living under an H2 section; the
// H2 sections are what the digesting pass folds into theme files. So an
// add must land the bullet INSIDE a section — above ## Themes (or the
// invariant breaks) and below the preamble (or nothing enumerates it).
//
// This replaces the pd-m1 proof, which asserted the retired model:
// "### " entries appended at EOF into a trailing "## Recent" section.

const rootConvention = `# Conventions

<!-- convention format … -->

## Go style

- gofmt clean before commit.

## Layout

- one package per concern.

## Themes

- naming — file and symbol naming rules → [naming](conventions/naming.md)
`

const newConvention = "- prefer errors.Is at boundaries\n"

// add applies the convention add-path with an optional target section.
func add(content, entry, section string) string {
	return string(insert.AppendEntry(
		[]byte(content), entry, cfgEntry.Convention, section,
	))
}

// T16: a named section receives the bullet at its top, newest-first,
// leaving the themes region and the other sections untouched.
func TestAdd_ConventionLandsInNamedSection(t *testing.T) {
	out := add(rootConvention, newConvention, "Layout")

	if !strings.Contains(out, "## Themes") ||
		!strings.Contains(out, "[naming](conventions/naming.md)") {
		t.Fatalf("Themes section/link destroyed by add; result:\n%s", out)
	}

	newIdx := strings.Index(out, "- prefer errors.Is at boundaries")
	layoutIdx := strings.Index(out, "## Layout")
	existingIdx := strings.Index(out, "- one package per concern.")
	themesIdx := strings.Index(out, "## Themes")
	if newIdx == -1 {
		t.Fatalf("new convention missing; result:\n%s", out)
	}
	if newIdx < layoutIdx {
		t.Errorf("bullet landed above its section heading (new=%d layout=%d)",
			newIdx, layoutIdx)
	}
	if newIdx > existingIdx {
		t.Errorf("bullet is not newest-first (new=%d existing=%d)",
			newIdx, existingIdx)
	}
	if newIdx > themesIdx {
		t.Errorf("bullet landed below ## Themes (new=%d themes=%d); it must "+
			"stay in the staging zone", newIdx, themesIdx)
	}
}

// T16: with no section named, the bullet goes to the first section —
// still inside staging, never the preamble.
func TestAdd_ConventionDefaultsToFirstSection(t *testing.T) {
	out := add(rootConvention, newConvention, "")

	newIdx := strings.Index(out, "- prefer errors.Is at boundaries")
	firstIdx := strings.Index(out, "## Go style")
	secondIdx := strings.Index(out, "## Layout")
	if newIdx == -1 {
		t.Fatalf("new convention missing; result:\n%s", out)
	}
	if newIdx < firstIdx || newIdx > secondIdx {
		t.Errorf("bullet did not land in the first section "+
			"(new=%d first=%d second=%d):\n%s",
			newIdx, firstIdx, secondIdx, out)
	}
}

// T16: an unknown section is created rather than silently retargeted,
// and it is placed above ## Themes so it lands in the staging zone.
func TestAdd_ConventionCreatesMissingSection(t *testing.T) {
	out := add(rootConvention, newConvention, "Error Handling")

	newSecIdx := strings.Index(out, "## Error Handling")
	themesIdx := strings.Index(out, "## Themes")
	newIdx := strings.Index(out, "- prefer errors.Is at boundaries")
	if newSecIdx == -1 {
		t.Fatalf("section not created; result:\n%s", out)
	}
	if newSecIdx > themesIdx {
		t.Errorf("created section sits below ## Themes (sec=%d themes=%d)",
			newSecIdx, themesIdx)
	}
	if newIdx < newSecIdx || newIdx > themesIdx {
		t.Errorf("bullet not inside the created section:\n%s", out)
	}
}

// T17 (C10): after a full fold the root has no sections left, only
// ## Themes. An add must still not land below it — the add-path creates
// a section above the themes region instead of appending at EOF.
//
// The CLI cannot reach this path: `ctx convention add` refuses an empty
// --section (see build.requireSection). This pins the fallback for
// direct AppendEntry callers, which lands in the deliberately
// unattractive "Unfiled" section rather than a comfortable catch-all.
func TestAdd_ConventionAfterFoldStaysAboveThemes(t *testing.T) {
	folded := "# Conventions\n\n<!-- guide -->\n\n" +
		"## Themes\n\n- naming — g → [naming](conventions/naming.md)\n"

	out := add(folded, newConvention, "")

	newIdx := strings.Index(out, "- prefer errors.Is at boundaries")
	themesIdx := strings.Index(out, "## Themes")
	if newIdx == -1 {
		t.Fatalf("new convention missing; result:\n%s", out)
	}
	if newIdx > themesIdx {
		t.Errorf("bullet landed below ## Themes on a folded root "+
			"(new=%d themes=%d); entry-below-themes would break:\n%s",
			newIdx, themesIdx, out)
	}
	if !strings.Contains(out, "## Unfiled") {
		t.Errorf("no section created to hold the bullet:\n%s", out)
	}
}
