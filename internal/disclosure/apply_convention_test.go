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

const conventionPreamble = "# Conventions\n\n<!-- guide -->\n\n"

// conventionStaging is three curated sections, each with the blank-line
// spacing a real root carries between them.
const conventionStaging = "## Error handling\n\nwrap with %w.\n\n" +
	"## Naming\n\nname files for their concern.\n\n" +
	"## Testing\n\ntable-driven, one case per row.\n\n"

// conventionRoot builds a migrated CONVENTIONS root: three staged
// sections above an existing one-bullet ## Themes.
func conventionRoot() string {
	return conventionPreamble + conventionStaging +
		"## Themes\n\n- style — house style → [style](conventions/style.md)\n"
}

// sec is a convention staged entry: title identity, no timestamp.
func sec(title string) disclosure.StagedEntry {
	return disclosure.StagedEntry{Title: title}
}

// writeConventionRoot writes content to <dir>/CONVENTIONS.md.
func writeConventionRoot(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "CONVENTIONS.md")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write root: %v", err)
	}
	return p
}

// M4/T14 (C5): folding a convention root moves whole sections into theme
// files, folds their gists, rewrites the root once, and leaves a root
// that still validates — with every moved byte conserved.
func TestApply_Convention(t *testing.T) {
	dir := t.TempDir()
	rootPath := writeConventionRoot(t, dir, conventionRoot())
	before := readFile(t, rootPath)

	plan := disclosure.Plan{
		Kind: "convention",
		Assignments: []disclosure.Assignment{
			{Theme: "errors", Slug: "errors", Gist: "how failures surface",
				Entries: []disclosure.StagedEntry{sec("Error handling")}},
			{Theme: "style", Slug: "style", Gist: "naming and test shape",
				Entries: []disclosure.StagedEntry{sec("Naming"), sec("Testing")}},
		},
	}

	res, err := disclosure.Apply(rootPath, plan, dir)
	if err != nil {
		t.Fatalf("Apply err = %v", err)
	}
	if res.Moved != 3 {
		t.Errorf("Moved = %d, want 3", res.Moved)
	}

	after := readFile(t, rootPath)
	root := disclosure.Parse(after, disclosure.KindConvention)

	// The folded root still satisfies every invariant.
	if vErr := disclosure.Validate(root); vErr != nil {
		t.Errorf("folded root fails Validate: %v", vErr)
	}
	// Staging is emptied of the moved sections — the root is bounded.
	if got := disclosure.StagedEntries(root); len(got) != 0 {
		t.Errorf("staging = %+v, want empty after folding all sections", got)
	}
	// The preamble survives verbatim.
	if !strings.HasPrefix(after, conventionPreamble) {
		t.Errorf("preamble not preserved; root=%q", after)
	}

	// Conservation: each section body lands in exactly one theme file and
	// is gone from the root.
	themeDir := filepath.Join(dir, "conventions")
	files, _ := filepath.Glob(filepath.Join(themeDir, "*.md"))
	for _, hdr := range []string{
		"## Error handling", "## Naming", "## Testing",
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
		if strings.Contains(root.Staging, hdr) {
			t.Errorf("%q still staged after the move", hdr)
		}
	}

	// Bodies travel with their headings.
	errorsFile := readFile(t, filepath.Join(themeDir, "errors.md"))
	if !strings.Contains(errorsFile, "wrap with %w.") {
		t.Errorf("errors.md missing its body: %q", errorsFile)
	}
	styleFile := readFile(t, filepath.Join(themeDir, "style.md"))
	for _, body := range []string{
		"name files for their concern.", "table-driven, one case per row.",
	} {
		if !strings.Contains(styleFile, body) {
			t.Errorf("style.md missing %q: %q", body, styleFile)
		}
	}

	// Gists folded, and the pre-existing bullet is untouched.
	if !strings.Contains(after, "how failures surface") {
		t.Errorf("errors gist not folded into the root: %q", after)
	}
	if !strings.Contains(after, "[errors](conventions/errors.md)") {
		t.Errorf("errors theme link missing: %q", after)
	}
	if before == after {
		t.Error("root unchanged after a non-empty fold")
	}
}

// M4/T15: a title-only plan flattens and addresses staged sections
// without a timestamp — the identity collapse the convention kind rests
// on. A title absent from staging is still rejected.
func TestApply_ConventionTitleIdentity(t *testing.T) {
	dir := t.TempDir()
	rootPath := writeConventionRoot(t, dir, conventionRoot())
	before := readFile(t, rootPath)

	_, err := disclosure.Apply(rootPath, disclosure.Plan{
		Kind: "convention",
		Assignments: []disclosure.Assignment{
			{Theme: "ghost", Slug: "ghost", Gist: "not there",
				Entries: []disclosure.StagedEntry{sec("No Such Section")}},
		},
	}, dir)
	if err == nil {
		t.Fatal("Apply err = nil, want a refusal for an unstaged title")
	}
	if after := readFile(t, rootPath); after != before {
		t.Error("root mutated on a refused plan")
	}
}
