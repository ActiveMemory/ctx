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

// Layout proof for specs/progressive-disclosure.md (plan pd-m1, T10).
//
// The design's load-bearing premise is that `ctx learning add` needs ZERO
// code change under progressive disclosure, because its insert anchor
// always lands the new entry in the staging zone — above the `## Themes`
// section — including when staging is empty and the anchor falls back to
// AfterHeader.
//
// These tests are the milestone's measurement gate: if they fail, the
// premise is wrong and the spec's Layout section must be revisited via
// /ctx-plan rather than patched downstream.

const themesSection = `## Themes

- hooks — hook mechanics, output channels, compliance → [hooks](learnings/hooks.md)
- output — write/ taxonomy and emission style → [output](learnings/output.md)
`

// rootPopulatedStaging: a Themes-bearing root that still has staged entries.
const rootPopulatedStaging = `# Learnings

<!--
UPDATE WHEN:
- Discover a gotcha worth recording
-->

## [2026-07-15-141726] gosec G101 has two independent triggers

**Context**: Removed the blanket exclusion.

---

` + themesSection

// rootEmptyStaging: a Themes-bearing root whose staging zone is empty —
// every entry has already been rolled out into its theme file. This is the
// steady state the design converges on, and the case that exercises the
// AfterHeader fallback.
const rootEmptyStaging = `# Learnings

<!--
UPDATE WHEN:
- Discover a gotcha worth recording
-->

` + themesSection

const newEntry = "## [2026-07-17-120000] A brand new learning\n\n" +
	"**Context**: freshly added.\n"

const newEntryHeader = "## [2026-07-17-120000]"

// assertLandsAboveThemes asserts the two things the premise requires: the
// Themes section survives the add, and the new entry sits above it (in the
// staging zone).
func assertLandsAboveThemes(t *testing.T, out string) {
	t.Helper()

	if !strings.Contains(out, "## Themes") {
		t.Fatalf("Themes section was DESTROYED by add; result:\n%s", out)
	}
	if !strings.Contains(out, "[hooks](learnings/hooks.md)") {
		t.Fatalf("theme gists/links were DESTROYED by add; result:\n%s", out)
	}

	entryIdx := strings.Index(out, newEntryHeader)
	themesIdx := strings.Index(out, "## Themes")
	if entryIdx == -1 {
		t.Fatalf("new entry missing from result:\n%s", out)
	}
	if entryIdx > themesIdx {
		t.Errorf(
			"new entry landed BELOW ## Themes (entry=%d themes=%d); "+
				"it must land in the staging zone above it",
			entryIdx, themesIdx,
		)
	}
}

// TestAdd_LearningLandsAboveThemes_PopulatedStaging covers the common case:
// staging holds entries, so the anchor finds a `## [` and inserts before it.
func TestAdd_LearningLandsAboveThemes_PopulatedStaging(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(rootPopulatedStaging), newEntry, cfgEntry.Learning, "",
	))
	assertLandsAboveThemes(t, out)
}

// TestAdd_LearningLandsAboveThemes_EmptyStaging covers the steady state:
// staging is empty, so the anchor finds no `## [` and falls back to
// AfterHeader. The premise says the entry still lands above ## Themes.
func TestAdd_LearningLandsAboveThemes_EmptyStaging(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(rootEmptyStaging), newEntry, cfgEntry.Learning, "",
	))
	assertLandsAboveThemes(t, out)
}
