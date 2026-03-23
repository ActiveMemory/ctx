//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

// finding describes a single knowledge file that exceeds its
// configured threshold.
type finding struct {
	// File is the context filename (e.g., DECISIONS.md).
	File string
	// Count is the actual entry or line count.
	Count int
	// Threshold is the configured maximum.
	Threshold int
	// Unit is the measurement unit ("entries" or "lines").
	Unit string
}
