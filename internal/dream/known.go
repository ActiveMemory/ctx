//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
)

// statusKnown reports whether s is a recognized ProposalStatus.
//
// Parameters:
//   - s: the status value to check
//
// Returns:
//   - bool: true when s is one of cfgDream.KnownStatuses
func statusKnown(s cfgDream.ProposalStatus) bool {
	for _, known := range cfgDream.KnownStatuses {
		if s == known {
			return true
		}
	}
	return false
}

// actionKnown reports whether a is a recognized ProposalAction.
//
// Parameters:
//   - a: the action value to check
//
// Returns:
//   - bool: true when a is one of cfgDream.KnownActions
func actionKnown(a cfgDream.ProposalAction) bool {
	for _, known := range cfgDream.KnownActions {
		if a == known {
			return true
		}
	}
	return false
}

// confidenceKnown reports whether c is a recognized Confidence level.
//
// Parameters:
//   - c: the confidence value to check
//
// Returns:
//   - bool: true when c is one of cfgDream.KnownConfidences
func confidenceKnown(c cfgDream.Confidence) bool {
	for _, known := range cfgDream.KnownConfidences {
		if c == known {
			return true
		}
	}
	return false
}
