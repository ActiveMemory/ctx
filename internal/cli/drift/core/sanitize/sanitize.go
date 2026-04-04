//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgDrift "github.com/ActiveMemory/ctx/internal/config/drift"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// FormatCheckName converts internal check identifiers to
// human-readable names.
//
// Parameters:
//   - name: Internal check identifier
//
// Returns:
//   - string: Human-readable description, or the original
//     name if unknown
func FormatCheckName(name cfgDrift.CheckName) string {
	switch name {
	case cfgDrift.CheckPathReferences:
		return desc.Text(text.DescKeyDriftCheckPathRefs)
	case cfgDrift.CheckStaleness:
		return desc.Text(text.DescKeyDriftCheckStaleness)
	case cfgDrift.CheckConstitution:
		return desc.Text(text.DescKeyDriftCheckConstitution)
	case cfgDrift.CheckRequiredFiles:
		return desc.Text(text.DescKeyDriftCheckRequired)
	case cfgDrift.CheckFileAge:
		return desc.Text(text.DescKeyDriftCheckFileAge)
	case cfgDrift.CheckTemplateHeaders:
		return desc.Text(
			text.DescKeyDriftCheckTemplateHeader,
		)
	default:
		return name
	}
}
