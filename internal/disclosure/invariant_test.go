//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// writeThemeFile creates dir/name with content and returns dir.
func writeThemeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if mkErr := os.MkdirAll(dir, 0o750); mkErr != nil {
		t.Fatalf("mkdir %s: %v", dir, mkErr)
	}
	if wErr := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); wErr != nil {
		t.Fatalf("write %s: %v", name, wErr)
	}
}

func themeWith(name, link string) disclosure.Root {
	return disclosure.Root{
		Kind:      disclosure.KindLearning,
		HasThemes: true,
		Themes:    []disclosure.Theme{{Name: name, Link: link}},
	}
}

// T07: gists <-> theme files must be 1:1.
func TestInvariant_Pairing(t *testing.T) {
	t.Run("1:1 passes", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings")
		writeThemeFile(t, dir, "hooks.md", "# hooks\n")
		root := themeWith("hooks", "learnings/hooks.md")
		if err := disclosure.CheckPairing(root, dir); err != nil {
			t.Errorf("CheckPairing = %v, want nil", err)
		}
	})

	t.Run("orphan theme file", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings")
		writeThemeFile(t, dir, "hooks.md", "# hooks\n")
		writeThemeFile(t, dir, "extra.md", "# extra\n") // no gist points here
		root := themeWith("hooks", "learnings/hooks.md")
		if err := disclosure.CheckPairing(root, dir); !errors.Is(err, errDisc.ErrOrphanThemeFile) {
			t.Errorf("CheckPairing = %v, want ErrOrphanThemeFile", err)
		}
	})

	t.Run("gist with no file", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings") // dir absent
		root := themeWith("hooks", "learnings/hooks.md")
		if err := disclosure.CheckPairing(root, dir); !errors.Is(err, errDisc.ErrMissingThemeFile) {
			t.Errorf("CheckPairing = %v, want ErrMissingThemeFile", err)
		}
	})

	t.Run("vacuous when un-migrated", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings") // absent
		root := disclosure.Root{Kind: disclosure.KindLearning}
		if err := disclosure.CheckPairing(root, dir); err != nil {
			t.Errorf("CheckPairing = %v, want nil (0<->0)", err)
		}
	})
}

// T08: every entry lives in exactly one place.
func TestInvariant_Uniqueness(t *testing.T) {
	const entry = "## [2026-07-15-120000] the entry\n\n**Context**: x.\n"

	t.Run("single occurrence passes", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings")
		writeThemeFile(t, dir, "hooks.md", entry)
		root := disclosure.Root{Kind: disclosure.KindLearning} // empty staging
		if err := disclosure.CheckUniqueness(root, dir); err != nil {
			t.Errorf("CheckUniqueness = %v, want nil", err)
		}
	})

	t.Run("dup across staging and theme file", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings")
		writeThemeFile(t, dir, "hooks.md", entry)
		root := disclosure.Root{Kind: disclosure.KindLearning, Staging: entry}
		if err := disclosure.CheckUniqueness(root, dir); !errors.Is(err, errDisc.ErrDuplicateEntry) {
			t.Errorf("CheckUniqueness = %v, want ErrDuplicateEntry", err)
		}
	})

	t.Run("dup across two theme files", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "learnings")
		writeThemeFile(t, dir, "hooks.md", entry)
		writeThemeFile(t, dir, "output.md", entry) // same entry in two files
		root := disclosure.Root{Kind: disclosure.KindLearning}
		if err := disclosure.CheckUniqueness(root, dir); !errors.Is(err, errDisc.ErrDuplicateEntry) {
			t.Errorf("CheckUniqueness = %v, want ErrDuplicateEntry", err)
		}
	})
}

// T09: every theme link must resolve.
func TestInvariant_Links(t *testing.T) {
	t.Run("resolving link passes", func(t *testing.T) {
		ctxDir := t.TempDir()
		writeThemeFile(t, filepath.Join(ctxDir, "learnings"), "hooks.md", "# hooks\n")
		root := themeWith("hooks", "learnings/hooks.md")
		if err := disclosure.CheckLinks(root, ctxDir); err != nil {
			t.Errorf("CheckLinks = %v, want nil", err)
		}
	})

	t.Run("broken link", func(t *testing.T) {
		ctxDir := t.TempDir() // nothing created
		root := themeWith("hooks", "learnings/hooks.md")
		if err := disclosure.CheckLinks(root, ctxDir); !errors.Is(err, errDisc.ErrBrokenThemeLink) {
			t.Errorf("CheckLinks = %v, want ErrBrokenThemeLink", err)
		}
	})
}
