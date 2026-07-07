//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// navMdEntry matches a quoted, docs-relative Markdown path in
// zensical.toml's nav, e.g. "recipes/spec-driven-development.md".
var navMdEntry = regexp.MustCompile(`"([^"]+\.md)"`)

// TestEveryDocPageIsReachableInNav asserts that every Markdown page
// under docs/ is referenced in zensical.toml's nav, so it is reachable
// from the site sidebar.
//
// This is a No-Broken-Windows guard (CONSTITUTION.md): a docs page that
// exists but is absent from the nav is silently unreachable — a user
// only finds it by guessing the URL. The gap is easy to introduce
// (add a recipe or CLI page, forget the nav entry) and invisible until
// someone audits by hand, which is how a batch of orphaned recipe and
// CLI pages accumulated. This test surfaces the next one in `go test`.
//
// Two directories are deliberately excluded:
//   - blog/: posts are rendered by the blog plugin from blog/index.md,
//     not listed individually in the manual nav.
//   - any includes/ dir: partials/snippets embedded into other pages,
//     never standalone nav entries.
func TestEveryDocPageIsReachableInNav(t *testing.T) {
	root := projectRoot(t)

	navBytes, err := os.ReadFile(filepath.Join(root, "zensical.toml"))
	if err != nil {
		t.Fatalf("read zensical.toml: %v", err)
	}
	nav := string(navBytes)

	docsDir := filepath.Join(root, "docs")
	var orphans []string
	checked := 0

	walkErr := filepath.WalkDir(docsDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}

			rel, relErr := filepath.Rel(docsDir, path)
			if relErr != nil {
				return relErr
			}
			rel = filepath.ToSlash(rel)

			// Excluded by design (see doc comment).
			if strings.HasPrefix(rel, "blog/") {
				return nil
			}
			if strings.HasPrefix(rel, "includes/") ||
				strings.Contains(rel, "/includes/") {
				return nil
			}

			checked++
			// Nav entries are quoted, docs-relative paths, e.g.
			// "recipes/spec-driven-development.md".
			if !strings.Contains(nav, `"`+rel+`"`) {
				orphans = append(orphans, rel)
			}
			return nil
		})
	if walkErr != nil {
		t.Fatalf("walk docs/: %v", walkErr)
	}

	// Sentinel: a wrong docsDir or an empty tree would make the walk a
	// no-op and the test pass vacuously. Require that it actually saw
	// pages, so the guard cannot silently stop guarding.
	if checked == 0 {
		t.Fatalf("no docs pages found under %s — test is not guarding anything",
			docsDir)
	}

	if len(orphans) > 0 {
		t.Errorf(
			"%d docs page(s) are absent from zensical.toml's nav and "+
				"thus unreachable from the site sidebar — add each under "+
				"the right nav group (or exclude it deliberately in this "+
				"test):\n  %s",
			len(orphans), strings.Join(orphans, "\n  "),
		)
	}

	// Reverse guard: every Markdown path named in the nav must resolve to
	// a real file. This catches the opposite drift — a page deleted or
	// renamed while its nav entry lingers, leaving a dead sidebar link.
	var deadLinks []string
	// Resolve through a rooted fs.FS: it rejects paths that escape docsDir
	// (no ".." traversal) by construction, so a stray nav entry cannot
	// stat an arbitrary file — and an invalid entry reads as a dead link.
	docsFS := os.DirFS(docsDir)
	for _, m := range navMdEntry.FindAllStringSubmatch(nav, -1) {
		if _, statErr := fs.Stat(docsFS, m[1]); statErr != nil {
			deadLinks = append(deadLinks, m[1])
		}
	}
	if len(deadLinks) > 0 {
		t.Errorf(
			"%d nav entr(y/ies) in zensical.toml point at a docs page that "+
				"does not exist — remove or fix each:\n  %s",
			len(deadLinks), strings.Join(deadLinks, "\n  "),
		)
	}
}
