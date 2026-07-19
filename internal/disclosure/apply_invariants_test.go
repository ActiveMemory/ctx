//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
)

// T08: after a successful Apply, all four disclosure invariants hold on
// the resulting root and its theme directory — Validate (structure),
// CheckPairing / CheckUniqueness (gists ↔ files 1:1, one place per
// entry), and CheckLinks (every link resolves).
func TestApplyInvariants(t *testing.T) {
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

	idA := idFor("2026-01-01-000000", "Alpha")
	idB := idFor("2026-01-02-000000", "Beta")
	plan := disclosure.Plan{
		Kind: "learning",
		Assignments: []disclosure.Assignment{
			{Theme: "context", Slug: "context", Gist: "context entries", Entries: []string{idA}},
			{Theme: "hooks", Slug: "hooks", Gist: "hook mechanics", Entries: []string{idB}},
		},
	}
	if _, err := disclosure.Apply(rootPath, plan, dir); err != nil {
		t.Fatalf("Apply err = %v", err)
	}

	root := disclosure.Parse(readFile(t, rootPath), disclosure.KindLearning)

	if err := disclosure.Validate(root); err != nil {
		t.Errorf("Validate = %v", err)
	}
	if err := disclosure.CheckPairing(root, learnDir); err != nil {
		t.Errorf("CheckPairing = %v", err)
	}
	if err := disclosure.CheckUniqueness(root, learnDir); err != nil {
		t.Errorf("CheckUniqueness = %v", err)
	}
	if err := disclosure.CheckLinks(root, dir); err != nil {
		t.Errorf("CheckLinks = %v", err)
	}
}
