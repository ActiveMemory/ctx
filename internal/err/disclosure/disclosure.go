//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
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

	// ErrNotAKnowledgeFile: the file handed to `ctx disclosure` is not a
	// canonical knowledge file (LEARNINGS/DECISIONS/CONVENTIONS). Wrap
	// with [NotAKnowledgeFile] to name the offending path.
	ErrNotAKnowledgeFile = entity.Sentinel(
		text.DescKeyErrDisclosureNotKnowledgeFileMsg,
	)
)

// Mover sentinels (milestone 3): the digesting pass that appends entry
// bodies to theme files, verifies byte-presence, and only then removes
// them from staging. All fail-loud with no auto-repair.
const (
	// ErrApplyNotEntryKind: the mover was handed a root whose kind it does
	// not digest — conventions (a later milestone) or an unknown kind. It
	// refuses rather than write to a guessed path.
	ErrApplyNotEntryKind = entity.Sentinel(
		text.DescKeyErrDisclosureApplyNotEntryKind,
	)

	// ErrEmptyAssignment: a plan assignment names no entries to move. A
	// theme with nothing to receive is a malformed plan, not a no-op.
	ErrEmptyAssignment = entity.Sentinel(
		text.DescKeyErrDisclosureEmptyAssignment,
	)

	// ErrEntryAssignedTwice: a plan assigns the same staged entry to more
	// than one theme — its destination is ambiguous.
	ErrEntryAssignedTwice = entity.Sentinel(
		text.DescKeyErrDisclosureEntryAssignedTwice,
	)

	// ErrEntryNotInStaging: a plan references an entry id absent from the
	// root's staging zone — the plan and the root disagree.
	ErrEntryNotInStaging = entity.Sentinel(
		text.DescKeyErrDisclosureEntryNotInStaging,
	)

	// ErrVerifyFailed: after appending an entry body to its theme file,
	// the body was not byte-present on re-read. The pass aborts with the
	// root untouched (append→verify→remove; never remove-then-append).
	ErrVerifyFailed = entity.Sentinel(
		text.DescKeyErrDisclosureVerifyFailed,
	)
)

// NotAKnowledgeFile wraps [ErrNotAKnowledgeFile] with the offending path
// and the expected filenames.
//
// Parameters:
//   - path: the file that was not a canonical knowledge file
//
// Returns:
//   - error: wrapping ErrNotAKnowledgeFile for errors.Is matches
func NotAKnowledgeFile(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDisclosureNotKnowledgeFile),
		ErrNotAKnowledgeFile,
		path,
	)
}
