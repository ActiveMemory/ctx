//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

import (
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
)

// allowed is the closed set of legal state transitions in the
// source-coverage state machine. Any (from, to) not present is
// rejected by [ValidTransition].
//
// Reading order matches the spec's table: every row's outgoing
// transitions appear left to right.
var allowed = map[transition]bool{
	{cfgKB.StateDiscovered, cfgKB.StateAdmitted}:          true,
	{cfgKB.StateDiscovered, cfgKB.StateSkipped}:           true,
	{cfgKB.StateAdmitted, cfgKB.StateHighlightsExtracted}: true,
	{cfgKB.StateAdmitted, cfgKB.StatePartiallyIngested}:   true,
	{cfgKB.StateAdmitted, cfgKB.StateTopicPageDrafted}:    true,
	{cfgKB.StateAdmitted, cfgKB.StateComprehensive}:       true,
	{cfgKB.StateHighlightsExtracted,
		cfgKB.StatePartiallyIngested}: true,
	{cfgKB.StateHighlightsExtracted,
		cfgKB.StateTopicPageDrafted}: true,
	{cfgKB.StateHighlightsExtracted,
		cfgKB.StateComprehensive}: true,
	{cfgKB.StatePartiallyIngested,
		cfgKB.StateTopicPageDrafted}: true,
	{cfgKB.StatePartiallyIngested,
		cfgKB.StateComprehensive}: true,
	{cfgKB.StateTopicPageDrafted,
		cfgKB.StateComprehensive}: true,
}

// ValidTransition reports whether advancing from state `from`
// to state `to` is allowed. Same-state transitions (idempotent
// "touch" of a row without state change) are always allowed.
//
// Parameters:
//   - from: current state (string from cfgKB.State*).
//   - to: next state (string from cfgKB.State*).
//
// Returns:
//   - bool: true when the transition is legal.
func ValidTransition(from, to string) bool {
	if from == to {
		return true
	}
	return allowed[transition{from: from, to: to}]
}
