//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// clampFraction bounds a budget fraction to the closed interval [0, 1];
// out-of-range configuration is coerced rather than trusted. Used by
// the budget-percentage getters so a hand-edited .ctxrc cannot push an
// allocation above the whole budget or below zero.
//
// Parameters:
//   - v: candidate fraction
//
// Returns:
//   - float64: v clamped to [0, 1]
func clampFraction(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
