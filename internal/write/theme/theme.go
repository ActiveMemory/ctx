//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package theme

import (
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/disclosure"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	"github.com/ActiveMemory/ctx/internal/i18n"
	"github.com/ActiveMemory/ctx/internal/slug"
)

// IsTarget reports whether a --section value addresses the themes region
// rather than a content section. "Themes" and "## Themes" both qualify,
// case-insensitively.
//
// Parameters:
//   - section: the caller-supplied --section value
//
// Returns:
//   - bool: true when the add targets the themes region
func IsTarget(section string) bool {
	s := i18n.Fold(strings.TrimSpace(section))
	bare := strings.TrimPrefix(
		cfgDisc.HeadingThemes, token.HeadingLevelTwoStart,
	)
	return s == i18n.Fold(bare) || s == i18n.Fold(cfgDisc.HeadingThemes)
}

// Apply declares a theme in a root: it folds the gist bullet into the
// root's ## Themes region (creating that region when absent) and creates
// the theme file the bullet links to.
//
// Both halves are required. A gist whose theme file does not exist fails
// the root's gist-to-file pairing invariant, so writing the bullet alone
// would leave a root that `ctx disclosure` refuses to touch.
//
// The content is "<name> — <gist>", split on the same em-dash separator
// the theme parser reads back, so what is written round-trips.
//
// Parameters:
//   - contextDir: the context directory theme files live under
//   - basename: the root's file name (e.g. "CONVENTIONS.md")
//   - existing: the root's current content
//   - content: the "<name> — <gist>" spec
//
// Returns:
//   - []byte: the rewritten root content
//   - error: a malformed spec, an unknown kind, or an IO error
func Apply(
	contextDir, basename, existing, content string,
) ([]byte, error) {
	name, gist, ok := splitSpec(content)
	if !ok {
		return nil, errAdd.ThemeSpec()
	}

	stem := slug.FromTitle(name)
	if stem == "" {
		return nil, errAdd.ThemeSpec()
	}

	rewritten, noun, addErr := disclosure.AddTheme(
		existing, basename, name, stem, gist,
	)
	if addErr != nil {
		return nil, addErr
	}

	// Create the theme file before returning the rewritten root: the
	// bullet is a pointer, and a pointer whose target does not exist is
	// the dangling-link state the pairing invariant exists to catch.
	if fileErr := ensureFile(
		contextDir, noun, stem, name,
	); fileErr != nil {
		return nil, fileErr
	}

	return []byte(rewritten), nil
}
