//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// TestWriteScopeAllowDeny exercises the write-scope allow/deny matrix,
// including the sanctioned specs/ promote crossing.
func TestWriteScopeAllowDeny(t *testing.T) {
	root := t.TempDir()

	cases := []struct {
		name    string
		target  string
		action  cfgDream.ProposalAction
		allowed bool
	}{
		{
			name:    "under dreams",
			target:  "dreams/20260607/p1.json",
			action:  cfgDream.ActionKeep,
			allowed: true,
		},
		{
			name:    "under ideas",
			target:  "ideas/note.md",
			action:  cfgDream.ActionArchive,
			allowed: true,
		},
		{
			name:    "specs allowed only on promote",
			target:  "specs/new.md",
			action:  cfgDream.ActionPromote,
			allowed: true,
		},
		{
			name:    "specs denied without promote",
			target:  "specs/new.md",
			action:  cfgDream.ActionKeep,
			allowed: false,
		},
		{
			name:    "tracked source dir denied",
			target:  "internal/x.go",
			action:  cfgDream.ActionKeep,
			allowed: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dec, err := dream.WriteScope(root, tc.target, tc.action)
			if err != nil {
				t.Fatalf("WriteScope error: %v", err)
			}
			if dec.Allowed != tc.allowed {
				t.Fatalf(
					"Allowed = %v, want %v (reason %q)",
					dec.Allowed, tc.allowed, dec.Reason,
				)
			}
			if !dec.Allowed && dec.Reason == "" {
				t.Fatal("refusal must carry a Reason")
			}
		})
	}
}

// gitInit initializes a throwaway repo at root with a .gitignore that
// ignores dreams/, and commits the .gitignore so specs/ is tracked.
func gitInit(t *testing.T, root string) {
	t.Helper()
	run := func(args ...string) {
		cmd := exec.Command("git", args...) //nolint:gosec // test fixture
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")
	gitignore := filepath.Join(root, ".gitignore")
	if writeErr := writeFixture(
		gitignore, "dreams/\nideas/\n",
	); writeErr != nil {
		t.Fatalf("write .gitignore: %v", writeErr)
	}
	run("add", ".gitignore")
	run("commit", "-m", "init")
}

// TestLeakIgnoredAllowed allows a write under a gitignored directory.
func TestLeakIgnoredAllowed(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)

	dec, err := dream.Leak(
		root, "dreams/20260607/p1.json", cfgDream.ActionKeep,
	)
	if err != nil {
		t.Fatalf("Leak error: %v", err)
	}
	if !dec.Allowed {
		t.Fatalf("gitignored path must be allowed (reason %q)", dec.Reason)
	}
}

// TestLeakTrackedRefused refuses a write to a tracked path.
func TestLeakTrackedRefused(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)

	dec, err := dream.Leak(
		root, ".gitignore", cfgDream.ActionKeep,
	)
	if err != nil {
		t.Fatalf("Leak error: %v", err)
	}
	if dec.Allowed {
		t.Fatal("tracked path must be refused")
	}
	if dec.Reason == "" {
		t.Fatal("refusal must carry a Reason")
	}
}

// TestLeakPromoteCrossing allows the specs/ promote crossing even
// though specs/ is tracked.
func TestLeakPromoteCrossing(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)

	dec, err := dream.Leak(
		root, "specs/new.md", cfgDream.ActionPromote,
	)
	if err != nil {
		t.Fatalf("Leak error: %v", err)
	}
	if !dec.Allowed {
		t.Fatalf("promote crossing must be allowed (reason %q)", dec.Reason)
	}
}
