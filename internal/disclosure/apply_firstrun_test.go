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

// T07: on an un-migrated root (no ## Themes), Apply creates the section
// below the remaining staging, leaves un-moved entries staged, and the
// result validates.
func TestApplyFirstRun(t *testing.T) {
	dir := t.TempDir()
	rootPath := writeRoot(t, dir, learnPreamble+stagingTwo) // no ## Themes

	idA := idFor("2026-01-01-000000", "Alpha")
	plan := disclosure.Plan{
		Kind: "learning",
		Assignments: []disclosure.Assignment{
			{Theme: "context", Slug: "context", Gist: "context entries", Entries: []string{idA}},
		},
	}

	if _, err := disclosure.Apply(rootPath, plan, dir); err != nil {
		t.Fatalf("Apply err = %v", err)
	}

	got := readFile(t, rootPath)
	if !strings.Contains(got, "## Themes") {
		t.Errorf("themes section not created:\n%q", got)
	}
	// Beta remains staged; Alpha moved out.
	if !strings.Contains(got, "## [2026-01-02-000000] Beta") {
		t.Errorf("un-moved Beta missing from staging:\n%q", got)
	}
	if strings.Contains(got, "## [2026-01-01-000000] Alpha") {
		t.Errorf("moved Alpha still in root:\n%q", got)
	}
	// The heading sits on its own line after a blank line.
	if !strings.Contains(got, "\n\n## Themes\n") {
		t.Errorf("themes boundary not blank-line separated:\n%q", got)
	}
	if err := disclosure.Validate(disclosure.Parse(got, disclosure.KindLearning)); err != nil {
		t.Errorf("post-first-run Validate = %v", err)
	}
}

// T07: an empty plan is a no-op — the root is byte-identical.
func TestApplyIdempotent(t *testing.T) {
	dir := t.TempDir()
	rootPath := writeRoot(t, dir, migratedRoot())
	before := readFile(t, rootPath)

	res, err := disclosure.Apply(rootPath, disclosure.Plan{Kind: "learning"}, dir)
	if err != nil {
		t.Fatalf("Apply err = %v", err)
	}
	if res.Moved != 0 {
		t.Errorf("Moved = %d, want 0", res.Moved)
	}
	if after := readFile(t, rootPath); after != before {
		t.Errorf("empty plan mutated the root")
	}
}
