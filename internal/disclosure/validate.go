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
	"github.com/ActiveMemory/ctx/internal/heading"
)

// Validate is the progressive-disclosure precondition (spec Guards §2).
// It refuses a structurally malformed root, so the pass never mutates
// content it does not understand — the failure mode behind the clobber
// bug. It is fail-loud with no auto-repair.
//
// Rules:
//   - zero or one "## Themes". Zero is the not-yet-migrated first-run
//     case and is valid; two or more is malformed (ErrMultipleThemes).
//   - no "## [" entry below "## Themes" (ErrEntryBelowThemes): entries
//     must stay in the staging zone above it.
//   - for entry kinds, a non-empty staging zone must parse into discrete
//     "## [" entries (ErrStagingUnparsable).
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

	entryBelow := r.HasThemes &&
		firstLinePrefixOffset(r.ThemesRaw, cfgDisc.EntryLinePrefix) != -1
	if entryBelow {
		return errDisc.ErrEntryBelowThemes
	}

	if r.Kind != KindConvention &&
		strings.TrimSpace(r.Staging) != "" &&
		len(heading.ParseEntryBlocks(r.Staging)) == 0 {
		return errDisc.ErrStagingUnparsable
	}

	return nil
}
