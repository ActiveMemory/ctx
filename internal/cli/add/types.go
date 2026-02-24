//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

// EntryParams contains all parameters needed to add an entry to a context file.
//
// Fields:
//   - Type: Entry type (decision, learning, task, convention)
//   - Content: Title or main content
//   - Section: Target section (for tasks)
//   - Priority: Priority level (for tasks)
//   - Context: Context field (for decisions/learnings)
//   - Rationale: Rationale (for decisions)
//   - Consequences: Consequences (for decisions)
//   - Lesson: Lesson (for learnings)
//   - Application: Application (for learnings)
type EntryParams struct {
	Type         string
	Content      string
	Section      string
	Priority     string
	Context      string
	Rationale    string
	Consequences string
	Lesson       string
	Application  string
}

// addConfig holds all flags for the add command.
//
// Fields:
//   - priority: Priority level for tasks (high, medium, low)
//   - section: Target section in TASKS.md
//   - fromFile: Read content from a file instead of argument
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
