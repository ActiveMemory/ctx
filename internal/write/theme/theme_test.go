//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package theme_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
	"github.com/ActiveMemory/ctx/internal/write/theme"
)

// M4/T16: only an explicit Themes section targets the themes region.
func TestIsTarget(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"Themes", true},
		{"themes", true},
		{"THEMES", true},
		{"## Themes", true},
		{"  Themes  ", true},
		{"", false},
		{"Naming", false},
		{"Theme", false},
		{"Themes and more", false},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			if got := theme.IsTarget(tc.in); got != tc.want {
				t.Errorf("IsTarget(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// M4/T16: declaring a theme writes the gist bullet AND creates the theme
// file it links to, so the gist-to-file pairing invariant holds at once.
func TestApply_WritesBulletAndFile(t *testing.T) {
	dir := t.TempDir()
	root := "# Conventions\n\n<!-- guide -->\n\n## Naming\n\n- use kebab-case.\n"

	out, err := theme.Apply(
		dir, "CONVENTIONS.md", root,
		"error handling — how failures surface",
	)
	if err != nil {
		t.Fatalf("Apply err = %v", err)
	}

	got := string(out)
	if !strings.Contains(got, "## Themes") {
		t.Errorf("themes region not created:\n%s", got)
	}
	if !strings.Contains(
		got, "[error handling](conventions/error-handling.md)",
	) {
		t.Errorf("theme link wrong or missing:\n%s", got)
	}
	if !strings.Contains(got, "how failures surface") {
		t.Errorf("gist missing:\n%s", got)
	}
	// Staging survives untouched — declaring a theme moves nothing.
	if !strings.Contains(got, "- use kebab-case.") {
		t.Errorf("existing staging lost:\n%s", got)
	}

	created := filepath.Join(dir, "conventions", "error-handling.md")
	body, readErr := os.ReadFile(created) //nolint:gosec // test temp path
	if readErr != nil {
		t.Fatalf("theme file not created: %v", readErr)
	}
	if !strings.HasPrefix(string(body), "# error handling") {
		t.Errorf("theme file header = %q", string(body))
	}

	// The rewritten root still validates.
	if vErr := disclosure.Validate(
		disclosure.Parse(got, disclosure.KindConvention),
	); vErr != nil {
		t.Errorf("root fails Validate after theme add: %v", vErr)
	}
}

// M4/T16: a spec without the name/gist separator is refused, and nothing
// is written — no bullet, no stray theme file.
func TestApply_RejectsMalformedSpec(t *testing.T) {
	dir := t.TempDir()
	root := "# Conventions\n\n## Naming\n\n- use kebab-case.\n"

	for _, spec := range []string{
		"just a name", "", " — gist only", "name — ",
	} {
		t.Run(spec, func(t *testing.T) {
			if _, err := theme.Apply(
				dir, "CONVENTIONS.md", root, spec,
			); err == nil {
				t.Fatalf("Apply(%q) err = nil, want a spec refusal", spec)
			}
			entries, _ := os.ReadDir(filepath.Join(dir, "conventions"))
			if len(entries) != 0 {
				t.Errorf("theme file created for a refused spec: %v", entries)
			}
		})
	}
}

// M4/T16: re-declaring a theme revises its gist without clobbering the
// existing theme file's contents.
func TestApply_ExistingThemeFilePreserved(t *testing.T) {
	dir := t.TempDir()
	themeDir := filepath.Join(dir, "conventions")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(themeDir, "errors.md")
	if err := os.WriteFile(
		path, []byte("# errors\n\nexisting body\n"), 0o600,
	); err != nil {
		t.Fatal(err)
	}

	root := "# Conventions\n\n## Themes\n\n- errors — old gist → " +
		"[errors](conventions/errors.md)\n"
	if _, err := theme.Apply(
		dir, "CONVENTIONS.md", root, "errors — revised gist",
	); err != nil {
		t.Fatalf("Apply err = %v", err)
	}

	body, readErr := os.ReadFile(path) //nolint:gosec // test temp path
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(body), "existing body") {
		t.Errorf("theme file clobbered: %q", string(body))
	}
}
