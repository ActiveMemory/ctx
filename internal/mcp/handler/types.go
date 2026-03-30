//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

// EntryOpts holds optional fields for entry creation.
//
// Fields:
//   - Priority: Priority label (high, medium, low)
//   - Context: Context field for decisions/learnings
//   - Rationale: Rationale field for decisions
//   - Consequence: Consequence field for decisions
//   - Lesson: Lesson field for learnings
//   - Application: Application field for learnings
type EntryOpts struct {
	Priority    string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
}
