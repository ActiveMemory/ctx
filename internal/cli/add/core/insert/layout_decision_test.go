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

// Layout proof for specs/progressive-disclosure.md (plan pd-m1, T11):
// DECISIONS. Decision() shares beforeFirstEntry with learnings, so the
// same premise applies — a new decision lands above ## Themes in both the
// populated- and empty-staging cases.

const rootDecisionPopulated = `# Decisions

<!-- DECISION FORMATS … -->

## [2026-07-15-000000] an existing decision

**Status**: Accepted

---

## Themes

- naming — singular CLI entity names → [naming](decisions/naming.md)
`

const rootDecisionEmpty = `# Decisions

<!-- DECISION FORMATS … -->

## Themes

- naming — singular CLI entity names → [naming](decisions/naming.md)
`

const newDecision = "## [2026-07-17-120000] A brand new decision\n\n" +
	"**Status**: Accepted\n"

func assertDecisionAboveThemes(t *testing.T, out string) {
	t.Helper()
	if !strings.Contains(out, "## Themes") ||
		!strings.Contains(out, "[naming](decisions/naming.md)") {
		t.Fatalf("Themes section/link destroyed by add; result:\n%s", out)
	}
	entryIdx := strings.Index(out, "## [2026-07-17-120000]")
	themesIdx := strings.Index(out, "## Themes")
	if entryIdx == -1 {
		t.Fatalf("new decision missing; result:\n%s", out)
	}
	if entryIdx > themesIdx {
		t.Errorf("decision landed below ## Themes (entry=%d themes=%d)",
			entryIdx, themesIdx)
	}
}

func TestAdd_DecisionLandsAboveThemes_PopulatedStaging(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(rootDecisionPopulated), newDecision, cfgEntry.Decision, "",
	))
	assertDecisionAboveThemes(t, out)
}

func TestAdd_DecisionLandsAboveThemes_EmptyStaging(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(rootDecisionEmpty), newDecision, cfgEntry.Decision, "",
	))
	assertDecisionAboveThemes(t, out)
}
