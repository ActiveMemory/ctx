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

// Layout proof for specs/progressive-disclosure.md (plan pd-m1, T12):
// CONVENTIONS. Unlike learnings/decisions, conventions append at EOF
// (AppendAtEnd) and their entries are `###` prose sections. So the design
// puts the staging zone in a trailing `## Recent` section (last in the
// file) — an EOF append lands the new section inside it, below `## Themes`.

const rootConvention = `# Conventions

<!-- convention format … -->

## Themes

- naming — file and symbol naming rules → [naming](conventions/naming.md)
- testing — colocation and integration split → [testing](conventions/testing.md)

## Recent

### an existing recent convention

Some prose.
`

const newConvention = "### a brand new convention\n\nMore prose.\n"

func TestAdd_ConventionLandsInRecent(t *testing.T) {
	out := string(insert.AppendEntry(
		[]byte(rootConvention), newConvention, cfgEntry.Convention, "",
	))

	// Themes section and its links must survive.
	if !strings.Contains(out, "## Themes") ||
		!strings.Contains(out, "[naming](conventions/naming.md)") {
		t.Fatalf("Themes section/link destroyed by add; result:\n%s", out)
	}

	newIdx := strings.Index(out, "### a brand new convention")
	recentIdx := strings.Index(out, "## Recent")
	themesIdx := strings.Index(out, "## Themes")
	if newIdx == -1 {
		t.Fatalf("new convention missing; result:\n%s", out)
	}
	// The new section must land inside ## Recent — i.e. after both the
	// Themes and the Recent headings.
	if newIdx < recentIdx {
		t.Errorf("convention landed ABOVE ## Recent (new=%d recent=%d); "+
			"it must land inside the trailing staging section", newIdx, recentIdx)
	}
	if newIdx < themesIdx {
		t.Errorf("convention landed above ## Themes (new=%d themes=%d)",
			newIdx, themesIdx)
	}
}
