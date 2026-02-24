//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import "github.com/ActiveMemory/ctx/internal/drift"

// formatCheckName converts internal check identifiers to human-readable names.
//
// Parameters:
//   - name: Internal check identifier
//     (e.g., "path_references", "staleness_check")
//
// Returns:
//   - string: Human-readable description of the check, or the original name
//     if unknown
func formatCheckName(name drift.CheckName) string {
	switch name {
	case drift.CheckPathReferences:
		return "Path references are valid"
	case drift.CheckStaleness:
		return "No staleness indicators"
	case drift.CheckConstitution:
		return "Constitution rules respected"
	case drift.CheckRequiredFiles:
		return "All required files present"
	case drift.CheckFileAge:
		return "No stale files by age"
	default:
		return string(name)
	}
}
