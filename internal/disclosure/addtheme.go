//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"path/filepath"

	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// AddTheme declares a theme in a root: it folds the gist bullet into the
// root's ## Themes region, creating that region when the root has none.
// Nothing moves out of staging — this is how a theme is named before any
// entries are digested into it.
//
// It takes and returns plain strings so callers outside this package do
// not have to construct a Root or an Assignment to declare a theme. The
// returned noun is the theme-file subdirectory the emitted gist links to,
// so the caller can create the file the link points at.
//
// Parameters:
//   - content: the root's current content
//   - basename: the root's file name (e.g. "CONVENTIONS.md")
//   - name: the theme's human-readable name (the bullet label)
//   - stem: the theme file's basename stem (the link target)
//   - gist: the "just enough" one-liner
//
// Returns:
//   - string: the rewritten root content
//   - string: the theme-file subdirectory for this kind
//   - error: ErrNotAKnowledgeFile, ErrApplyNotEntryKind, or nil
func AddTheme(
	content, basename, name, stem, gist string,
) (string, string, error) {
	kind, known := KindFor(filepath.Base(basename))
	if !known {
		return "", "", errDisc.NotAKnowledgeFile(basename)
	}
	noun, hasDir := ThemeDir(kind)
	if !hasDir {
		return "", "", errDisc.ErrApplyNotEntryKind
	}

	root := Parse(content, kind)
	assignment := Assignment{Theme: name, Slug: stem, Gist: gist}
	rewritten := rewriteRoot(
		root, root.Staging,
		Plan{Assignments: []Assignment{assignment}},
		noun,
	)
	return rewritten, noun, nil
}
