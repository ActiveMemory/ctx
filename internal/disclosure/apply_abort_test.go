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

// T06(b): when a theme append fails, Apply aborts before the root rewrite,
// so the root is byte-identical — the append→verify→remove ordering means
// the root write is the last syscall and never runs on a failed pass.
// The earlier theme append is allowed to remain (duplication is
// recoverable; loss is not — spec Guards §1,§3).
func TestApplyAbort_RootUntouched(t *testing.T) {
	dir := t.TempDir()
	// Make learnings/hooks.md a directory so the second append fails.
	if err := os.MkdirAll(filepath.Join(dir, "learnings", "hooks.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	rootPath := writeRoot(t, dir, migratedRoot())
	before := readFile(t, rootPath)

	plan := disclosure.Plan{
		Kind: "learning",
		Assignments: []disclosure.Assignment{
			// context: a fresh file, its append succeeds.
			{Theme: "context", Slug: "context", Gist: "g",
				Entries: []disclosure.StagedEntry{ent("2026-01-01-000000", "Alpha")}},
			// hooks: hooks.md is a directory, so this append fails.
			{Theme: "hooks", Slug: "hooks", Gist: "g",
				Entries: []disclosure.StagedEntry{ent("2026-01-02-000000", "Beta")}},
		},
	}

	if _, err := disclosure.Apply(rootPath, plan, dir); err == nil {
		t.Fatal("Apply err = nil, want a theme-append failure")
	}
	if after := readFile(t, rootPath); after != before {
		t.Errorf("root mutated on abort:\nbefore=%q\nafter =%q", before, after)
	}
	// The earlier (context) append did land — duplication over loss.
	if _, err := os.Stat(filepath.Join(dir, "learnings", "context.md")); err != nil {
		t.Errorf("earlier theme append missing: %v", err)
	}
}
