//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "github.com/ActiveMemory/ctx/internal/config/stats"

// percent computes the percentage of used relative to total.
//
// Returns 0 when total is zero to avoid division by zero.
//
// Parameters:
//   - used: Numerator value
//   - total: Denominator value
//
// Returns:
//   - float64: Percentage (0-100)
func percent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * stats.PercentMultiplier
}
