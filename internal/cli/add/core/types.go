//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// EntryParams contains all parameters needed to add an entry to a context file.
type EntryParams struct {
	Type        string
	Content     string
	Section     string
	Priority    string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
	ContextDir  string
}

// Config holds all flags for the add command.
//
// Fields:
//   - Priority: Priority level for tasks (high, medium, low)
//   - Section: Target section in TASKS.md
//   - FromFile: Read content from a file instead of argument
//   - Context: Context field for decisions/learnings
//   - Rationale: Rationale field for decisions
//   - Consequence: Consequence field for decisions
//   - Lesson: Lesson field for learnings
//   - Application: Application field for learnings
type Config struct {
	Priority    string
	Section     string
	FromFile    string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
}
