//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// ImportAction describes what will happen to a given file.
type ImportAction int

const (
	ActionNew        ImportAction = iota // file does not exist yet
	ActionRegenerate                     // file exists and will be rewritten
	ActionSkip                           // file exists and will be left alone
	ActionLocked                         // file is locked: never overwritten
)

// ImportOpts holds all flag values for the import command.
//
// Fields:
//   - All, AllProjects, Regenerate, Yes, DryRun: Boolean flags
//   - KeepFrontmatter: Preserve enriched YAML frontmatter during regen
type ImportOpts struct {
	All, AllProjects, Regenerate, Yes, DryRun bool
	KeepFrontmatter                           bool
}

// DiscardFrontmatter reports whether frontmatter should be discarded
// during regeneration.
//
// Returns:
//   - bool: True if frontmatter should be discarded
func (o ImportOpts) DiscardFrontmatter() bool {
	return !o.KeepFrontmatter
}

// FileAction describes the planned action for a single import file (one part
// of one session).
//
// Fields:
//   - Session: Source session metadata
//   - Filename: Output markdown filename
//   - Path: Full output path
//   - Part: Part number (1-based) for split sessions
//   - TotalParts: Total parts for this session
//   - StartIdx: First message index in this part
//   - EndIdx: Last message index (exclusive)
//   - Action: Planned action (new, regenerate, skip, locked)
//   - Messages: Messages belonging to this part
//   - Slug: URL-safe session identifier
//   - Title: Human-readable session title
//   - BaseName: Filename without extension
type FileAction struct {
	Session    *Session
	Filename   string
	Path       string
	Part       int
	TotalParts int
	StartIdx   int
	EndIdx     int
	Action     ImportAction
	Messages   []Message
	Slug       string
	Title      string
	BaseName   string
}

// ImportPlan is the result of plan.Import: a list of per-file actions plus
// aggregate counters and any renames that need to happen first.
//
// Fields:
//   - Actions: Per-file planned actions
//   - NewCount: Files that will be created
//   - RegenCount: Files that will be regenerated
//   - SkipCount: Files that will be skipped
//   - LockedCount: Files that are locked
//   - RenameOps: Dedup renames to execute before import
type ImportPlan struct {
	Actions     []FileAction
	NewCount    int
	RegenCount  int
	SkipCount   int
	LockedCount int
	RenameOps   []RenameOp
}

// RenameOp describes a dedup rename (old slug → new slug).
//
// Fields:
//   - OldBase: Original filename base
//   - NewBase: Deduplicated filename base
//   - NumParts: Number of parts to rename
type RenameOp struct {
	OldBase  string
	NewBase  string
	NumParts int
}

// ImportResult tracks per-type counts for memory import operations.
//
// Fields:
//   - Conventions: convention entries imported
//   - Decisions: decision entries imported
//   - Learnings: learning entries imported
//   - Tasks: task entries imported
//   - Skipped: entries skipped (unclassified)
//   - Dupes: duplicate entries skipped
type ImportResult struct {
	Conventions int
	Decisions   int
	Learnings   int
	Tasks       int
	Skipped     int
	Dupes       int
}

// Total returns the number of entries actually imported (excludes
// skips and duplicates).
//
// Returns:
//   - int: sum of conventions, decisions, learnings, and tasks
func (r ImportResult) Total() int {
	return r.Conventions + r.Decisions + r.Learnings + r.Tasks
}
