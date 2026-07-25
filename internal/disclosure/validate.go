//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// Validate is the progressive-disclosure precondition (spec Guards §2).
// It refuses a structurally malformed root, so the pass never mutates
// content it does not understand — the failure mode behind the clobber
// bug. It is fail-loud with no auto-repair.
//
// Rules:
//   - zero or one "## Themes". Zero is the not-yet-migrated first-run
//     case and is valid; two or more is malformed (ErrMultipleThemes).
//   - no entry heading below "## Themes" (ErrEntryBelowThemes): entries
//     must stay in the staging zone above it. The heading that opens an
//     entry is per-kind ([EntryPrefix]).
//   - a non-empty staging zone must enumerate into discrete entries
//     (ErrStagingUnparsable), for every kind.
//   - for conventions, no two staged sections share a title
//     (ErrDuplicateStagedTitle): with no timestamp, the title is the
//     whole identity, so a duplicate makes a plan entry unaddressable.
//
// Invariants that are vacuously true on an un-migrated root (no themes,
// no theme files) need no special-casing here.
//
// Parameters:
//   - r: a parsed root (from Parse)
//
// Returns:
//   - error: one of the disclosure sentinels, or nil when well-formed
func Validate(r Root) error {
	if len(headingLineOffsets(r.Reconstruct(), cfgDisc.HeadingThemes)) > 1 {
		return errDisc.ErrMultipleThemes
	}

	if r.HasThemes && entryBelowThemes(r.ThemesRaw, r.Kind) {
		return errDisc.ErrEntryBelowThemes
	}

	blocks := stagedBlocks(r.Staging, r.Kind)
	if strings.TrimSpace(r.Staging) != "" && len(blocks) == 0 {
		return errDisc.ErrStagingUnparsable
	}

	// Title-only identity is unique to conventions; entry kinds carry a
	// timestamp and may legitimately repeat a title, so this rule stays
	// scoped rather than becoming a general uniqueness check.
	if r.Kind == KindConvention {
		seen := make(map[string]bool, len(blocks))
		for _, b := range blocks {
			if seen[b.Title] {
				return errDisc.ErrDuplicateStagedTitle
			}
			seen[b.Title] = true
		}
	}

	return nil
}
