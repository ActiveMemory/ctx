//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/drift"
)

// FormatCheckName converts internal check identifiers to human-readable names.
//
// Parameters:
//   - name: Internal check identifier
//     (e.g., "path_references", "staleness_check")
//
// Returns:
//   - string: Human-readable description of the check, or the original name
//     if unknown
func FormatCheckName(name drift.CheckName) string {
	switch name {
	case drift.CheckPathReferences:
		return desc.Text(text.DescKeyDriftCheckPathRefs)
	case drift.CheckStaleness:
		return desc.Text(text.DescKeyDriftCheckStaleness)
	case drift.CheckConstitution:
		return desc.Text(text.DescKeyDriftCheckConstitution)
	case drift.CheckRequiredFiles:
		return desc.Text(text.DescKeyDriftCheckRequired)
	case drift.CheckFileAge:
		return desc.Text(text.DescKeyDriftCheckFileAge)
	case drift.CheckTemplateHeaders:
		return desc.Text(text.DescKeyDriftCheckTemplateHeader)
	default:
		return string(name)
	}
}
