//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

// Subdir is the per-session handover artifact subdirectory
// under `.context/`. Files inside are named
// `<TS>-<slug>.md` so concurrent agent runs never overwrite
// each other.
const Subdir = "handovers"

// Section header constants used when composing the handover
// markdown body. These are structural identifiers, not
// localizable prose: the read side matches them exactly.
const (
	// SectionSummary headers the past-tense session summary.
	SectionSummary = "## Summary"
	// SectionNext headers the future-tense first action for
	// the next session.
	SectionNext = "## Next session"
	// SectionHighlights headers the optional highlights list.
	SectionHighlights = "## Highlights"
	// SectionOpenQuestions headers the optional
	// open-questions list.
	SectionOpenQuestions = "## Open questions"
	// SectionFoldedCloseouts headers the auto-generated list
	// of closeouts folded into the new handover.
	SectionFoldedCloseouts = "## Folded closeouts"
)

// BranchDetached is the literal branch label written into
// handovers when git HEAD is not on a symbolic ref.
const BranchDetached = "detached"

// DefaultSlug is the fallback slug for titles that reduce to
// an empty string after kebab-case normalisation.
const DefaultSlug = "handover"

// FoldEntryPrefix is the marker that opens each `## Folded
// closeouts` list item.
const FoldEntryPrefix = "- `"

// FoldEntryModePrefix opens the metadata section in a folded
// entry by closing the filename backtick and introducing the
// mode field with a colon.
const FoldEntryModePrefix = "`: mode="

// FoldEntryPassModePrefix opens the optional pass-mode field
// in a folded entry's metadata section.
//
//nolint:gosec // G101: literal markdown delimiter, not a credential
const FoldEntryPassModePrefix = ", pass-mode="

// FoldEntryGeneratedAtPrefix opens the generated-at timestamp
// in a folded entry's metadata section.
const FoldEntryGeneratedAtPrefix = ", generated-at="
