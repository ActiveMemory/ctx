//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

// Placeholder values rejected from required body flags on
// ctx decision add / ctx learning add. Matching is exact,
// case-insensitive, after whitespace trimming. Substring
// matches are allowed (a real sentence may legitimately
// contain the word "TBD" inside longer prose).
const (
	// PlaceholderTBD matches the canonical "to be defined" marker.
	PlaceholderTBD = "tbd"
	// PlaceholderNA matches the "not applicable" marker (slash form).
	PlaceholderNA = "n/a"
	// PlaceholderNAShort matches the "not applicable" marker (compact).
	PlaceholderNAShort = "na"
	// PlaceholderNone matches a bare "none" answer.
	PlaceholderNone = "none"
	// PlaceholderSeeChat matches deferral to chat transcript.
	PlaceholderSeeChat = "see chat"
	// PlaceholderSeeAbove matches deferral to an earlier passage.
	PlaceholderSeeAbove = "see above"
	// PlaceholderSeeBelow matches deferral to a later passage.
	PlaceholderSeeBelow = "see below"
	// PlaceholderPending matches a "will be filled in later" marker.
	PlaceholderPending = "pending"
	// PlaceholderToBeDone matches the long form of the TBD marker.
	PlaceholderToBeDone = "to be done"
)

// Placeholders is the closed set used by the rejection check.
// Keys are stored lowercase; callers must lowercase the trimmed
// input before lookup.
var Placeholders = map[string]struct{}{
	PlaceholderTBD:      {},
	PlaceholderNA:       {},
	PlaceholderNAShort:  {},
	PlaceholderNone:     {},
	PlaceholderSeeChat:  {},
	PlaceholderSeeAbove: {},
	PlaceholderSeeBelow: {},
	PlaceholderPending:  {},
	PlaceholderToBeDone: {},
}
