//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
)

// T09: conservation — every entry staged before the pass is, after it,
// either still staged or present in exactly one theme file. Zero loss,
// zero duplication.
func TestApplyConservation(t *testing.T) {
	dir := t.TempDir()
	learnDir := filepath.Join(dir, "learnings")
	if err := os.MkdirAll(learnDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(learnDir, "hooks.md"),
		[]byte("# hooks\n\nexisting\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rootPath := writeRoot(t, dir, migratedRoot())

	before := disclosure.Parse(readFile(t, rootPath), disclosure.KindLearning)
	beforeStaged := stagedHeaders(before)
	if len(beforeStaged) != 2 {
		t.Fatalf("precondition: %d staged, want 2", len(beforeStaged))
	}

	idA := idFor("2026-01-01-000000", "Alpha")
	idB := idFor("2026-01-02-000000", "Beta")
	plan := disclosure.Plan{
		Kind: "learning",
		Assignments: []disclosure.Assignment{
			{Theme: "context", Slug: "context", Gist: "g1", Entries: []string{idA}},
			{Theme: "hooks", Slug: "hooks", Gist: "g2", Entries: []string{idB}},
		},
	}
	res, err := disclosure.Apply(rootPath, plan, dir)
	if err != nil {
		t.Fatalf("Apply err = %v", err)
	}

	after := disclosure.Parse(readFile(t, rootPath), disclosure.KindLearning)
	afterStaged := stagedHeaders(after)

	// staging_before == moved + staging_after (count).
	if res.Moved+len(afterStaged) != len(beforeStaged) {
		t.Errorf("moved(%d) + after(%d) != before(%d)",
			res.Moved, len(afterStaged), len(beforeStaged))
	}

	// Each moved entry appears in exactly one theme file, and nowhere in
	// staging.
	files, _ := filepath.Glob(filepath.Join(learnDir, "*.md"))
	for _, hdr := range []string{
		"## [2026-01-01-000000] Alpha",
		"## [2026-01-02-000000] Beta",
	} {
		got := 0
		for _, f := range files {
			if strings.Contains(readFile(t, f), hdr) {
				got++
			}
		}
		if got != 1 {
			t.Errorf("%q present in %d theme files, want exactly 1", hdr, got)
		}
		if strings.Contains(after.Staging, hdr) {
			t.Errorf("%q still in staging after move", hdr)
		}
	}
}

// stagedHeaders returns the staged entries' identities for a root.
func stagedHeaders(r disclosure.Root) []string {
	var ids []string
	for _, e := range disclosure.StagedEntries(r) {
		ids = append(ids, e.Timestamp+"|"+e.Title)
	}
	return ids
}
