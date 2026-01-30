//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

// addConfig holds all flags for the add command.
//
// Fields:
//   - priority: Priority level for tasks (high, medium, low)
//   - section: Target section in TASKS.md
//   - fromFile: Read content from file instead of argument
//   - context: Context field for decisions/learnings
//   - rationale: Rationale field for decisions
//   - consequences: Consequences field for decisions
//   - lesson: Lesson field for learnings
//   - application: Application field for learnings
type addConfig struct {
	priority     string
	section      string
	fromFile     string
	context      string
	rationale    string
	consequences string
	lesson       string
	application  string
}
