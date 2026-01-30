//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

// ContextUpdate represents a parsed context update command.
//
// Extracted from <context-update> XML tags in the input stream.
//
// Fields:
//   - Type: Update type (task, decision, learning, convention, complete)
//   - Content: The entry text (title) or search query for complete
//   - Context: Context field for learnings/decisions (what prompted this)
//   - Lesson: Lesson field for learnings (the key insight)
//   - Application: Application field for learnings (how to apply going forward)
//   - Rationale: Rationale field for decisions (why this choice)
//   - Consequences: Consequences field for decisions (what changes as a result)
type ContextUpdate struct {
	Type         string
	Content      string
	Context      string
	Lesson       string
	Application  string
	Rationale    string
	Consequences string
}
