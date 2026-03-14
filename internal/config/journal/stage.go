//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Journal processing stage names.
const (
	// StageExported marks a journal entry as exported from Claude Code.
	StageExported = "exported"
	// StageEnriched marks a journal entry as enriched with metadata.
	StageEnriched = "enriched"
	// StageNormalized marks a journal entry as normalized for rendering.
	StageNormalized = "normalized"
	// StageFencesVerified marks a journal entry as having verified code fences.
	StageFencesVerified = "fences_verified"
	// StageLocked marks a journal entry as locked (read-only).
	StageLocked = "locked"
)
