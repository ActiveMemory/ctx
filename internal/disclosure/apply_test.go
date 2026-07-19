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
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

const learnPreamble = "# Learnings\n\n<!-- guide -->\n\n"

// migratedRoot builds a LEARNINGS root with the two-entry staging and one
// existing "hooks" theme.
func migratedRoot() string {
	return learnPreamble + stagingTwo +
		"## Themes\n\n- hooks — old → [hooks](learnings/hooks.md)\n"
}

// writeRoot writes content to <dir>/LEARNINGS.md and returns the path.
func writeRoot(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "LEARNINGS.md")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}
	return p
}

// T05: Apply moves entries to their theme files (creating and appending),
// folds gists, and rewrites the root with the moved entries gone.
func TestApply(t *testing.T) {
	dir := t.TempDir()
	// Pre-existing hooks theme file, so the append path is exercised.
	if err := os.MkdirAll(filepath.Join(dir, "learnings"), 0o755); err != nil {
		t.Fatal(err)
	}
	hooksPath := filepath.Join(dir, "learnings", "hooks.md")
	if err := os.WriteFile(hooksPath, []byte("# hooks\n\nexisting\n"), 0o644); err != nil {
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

	res, err := disclosure.Apply(rootPath, plan, dir)
	if err != nil {
		t.Fatalf("Apply err = %v", err)
	}
	if res.Moved != 2 {
		t.Errorf("Moved = %d, want 2", res.Moved)
	}
	if len(res.Themes) != 2 || res.Themes[0] != "context" || res.Themes[1] != "hooks" {
		t.Errorf("Themes = %v, want [context hooks]", res.Themes)
	}

	// context.md created with Alpha's body.
	ctxBody := readFile(t, filepath.Join(dir, "learnings", "context.md"))
	if !strings.Contains(ctxBody, spanAlpha) {
		t.Errorf("context.md missing Alpha span:\n%q", ctxBody)
	}
	// hooks.md keeps its original content and gains Beta.
	hooksBody := readFile(t, hooksPath)
	if !strings.Contains(hooksBody, "existing") || !strings.Contains(hooksBody, spanBeta) {
		t.Errorf("hooks.md = %q", hooksBody)
	}
	// Root: staging emptied, gists folded, still valid.
	newRoot := readFile(t, rootPath)
	if strings.Contains(newRoot, "## [2026-01-01-000000] Alpha") {
		t.Errorf("Alpha still in root:\n%q", newRoot)
	}
	if !strings.Contains(newRoot, "- hooks — hook mechanics → [hooks](learnings/hooks.md)") {
		t.Errorf("hooks gist not updated:\n%q", newRoot)
	}
	if !strings.Contains(newRoot, "- context — context entries → [context](learnings/context.md)") {
		t.Errorf("context gist not added:\n%q", newRoot)
	}
	if vErr := disclosure.Validate(disclosure.Parse(newRoot, disclosure.KindLearning)); vErr != nil {
		t.Errorf("post-apply Validate = %v", vErr)
	}
}

// T05: Apply refuses a non-knowledge file, an un-digestible kind, and a
// malformed root — each before any write.
func TestApply_Refusals(t *testing.T) {
	dir := t.TempDir()

	t.Run("non-knowledge file", func(t *testing.T) {
		p := filepath.Join(dir, "README.md")
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := disclosure.Apply(p, disclosure.Plan{}, dir)
		if !errors.Is(err, errDisc.ErrNotAKnowledgeFile) {
			t.Errorf("err = %v, want ErrNotAKnowledgeFile", err)
		}
	})

	t.Run("convention kind refused", func(t *testing.T) {
		p := filepath.Join(dir, "CONVENTIONS.md")
		if err := os.WriteFile(p, []byte("# Conventions\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := disclosure.Apply(p, disclosure.Plan{}, dir)
		if err != errDisc.ErrApplyNotEntryKind {
			t.Errorf("err = %v, want ErrApplyNotEntryKind", err)
		}
	})

	t.Run("malformed root", func(t *testing.T) {
		bad := learnPreamble + "## Themes\n\n## Themes\n"
		p := writeRoot(t, dir, bad)
		before := readFile(t, p)
		_, err := disclosure.Apply(p, disclosure.Plan{
			Assignments: []disclosure.Assignment{
				{Theme: "t", Slug: "t", Gist: "g", Entries: []string{"x"}},
			},
		}, dir)
		if err != errDisc.ErrMultipleThemes {
			t.Errorf("err = %v, want ErrMultipleThemes", err)
		}
		if after := readFile(t, p); after != before {
			t.Errorf("root mutated on refusal")
		}
	})
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
