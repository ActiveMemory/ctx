//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// EntryParams contains all parameters needed to add an entry to a context file.
//
// Fields:
//   - Type: Entry type (decision, learning, convention, task)
//   - Content: Main entry text
//   - Section: Target section within the file
//   - Priority: Priority label (high, medium, low)
//   - SessionID: AI session identifier for task provenance
//   - Branch: Git branch name for task provenance
//   - Commit: Git commit hash for task provenance
//   - Context: Context field for decisions/learnings
//   - Rationale: Rationale field for decisions
//   - Consequence: Consequence field for decisions
//   - Lesson: Lesson field for learnings
//   - Application: Application field for learnings
//   - ContextDir: Path to the context directory
type EntryParams struct {
	Type        string
	Content     string
	Section     string
	Priority    string
	SessionID   string
	Branch      string
	Commit      string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
	ContextDir  string
}

// AddConfig holds all flags for the add command.
//
// Fields:
//   - Priority: Priority label flag
//   - Section: Target section flag
//   - FromFile: Path to read content from a file
//   - SessionID: AI session identifier for task provenance
//   - Branch: Git branch name for task provenance
//   - Commit: Git commit hash for task provenance
//   - Context: Context flag for decisions/learnings
//   - Rationale: Rationale flag for decisions
//   - Consequence: Consequence flag for decisions
//   - Lesson: Lesson flag for learnings
//   - Application: Application flag for learnings
//   - Share: Also publish to the shared hub
type AddConfig struct {
	Priority    string
	Section     string
	FromFile    string
	SessionID   string
	Branch      string
	Commit      string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
}

// EntryOpts holds optional fields for entry creation via MCP.
//
// Fields:
//   - Priority: Priority label (high, medium, low)
//   - Section: Target section for tasks (required for tasks)
//   - SessionID: AI session identifier for provenance
//   - Branch: Git branch name for provenance
//   - Commit: Git commit hash for provenance
//   - Context: Context field for decisions/learnings
//   - Rationale: Rationale field for decisions
//   - Consequence: Consequence field for decisions
//   - Lesson: Lesson field for learnings
//   - Application: Application field for learnings
type EntryOpts struct {
	Priority    string
	Section     string
	SessionID   string
	Branch      string
	Commit      string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
	Share       bool
}
