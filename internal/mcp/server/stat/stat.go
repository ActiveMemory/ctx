//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stat

// TotalAdds sums all entry add counts.
//
// Parameters:
//   - m: map of entry type to add count
//
// Returns:
//   - int: total adds across all types
func TotalAdds(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}
