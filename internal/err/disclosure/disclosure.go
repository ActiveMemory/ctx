//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrMultipleThemes: the root has more than one "## Themes" section.
	// Validate refuses rather than guess which is authoritative.
	ErrMultipleThemes = entity.Sentinel(
		text.DescKeyErrDisclosureMultipleThemes,
	)

	// ErrEntryBelowThemes: a "## [" entry sits below "## Themes". Entries
	// must stay in the staging zone above it, or the pass cannot find them.
	ErrEntryBelowThemes = entity.Sentinel(
		text.DescKeyErrDisclosureEntryBelowThemes,
	)

	// ErrStagingUnparsable: the staging zone did not parse into discrete
	// entries. Refuse rather than regenerate from "what was recognized".
	ErrStagingUnparsable = entity.Sentinel(
		text.DescKeyErrDisclosureStagingUnparsable,
	)

	// ErrOrphanThemeFile: a theme file exists with no matching gist in
	// the root — the root no longer reaches it (link-graph orphan).
	ErrOrphanThemeFile = entity.Sentinel(
		text.DescKeyErrDisclosureOrphanThemeFile,
	)

	// ErrMissingThemeFile: a root gist links to a theme file that does
	// not exist.
	ErrMissingThemeFile = entity.Sentinel(
		text.DescKeyErrDisclosureMissingThemeFile,
	)

	// ErrDuplicateEntry: an entry appears in more than one place (staging
	// and a theme file, or two theme files) — the "exactly one place"
	// invariant is broken.
	ErrDuplicateEntry = entity.Sentinel(
		text.DescKeyErrDisclosureDuplicateEntry,
	)

	// ErrBrokenThemeLink: a theme link in the root points at a path that
	// does not resolve on disk.
	ErrBrokenThemeLink = entity.Sentinel(
		text.DescKeyErrDisclosureBrokenThemeLink,
	)
)
