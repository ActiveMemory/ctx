//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
)

// ExportAction describes what will happen to a given file.
type ExportAction int

const (
	ActionNew        ExportAction = iota // file does not exist yet
	ActionRegenerate                     // file exists and will be rewritten
	ActionSkip                           // file exists and will be left alone
	ActionLocked                         // file is locked — never overwritten
)

// ExportOpts holds all flag values for the export command.
type ExportOpts struct {
	All, AllProjects, Force, Regenerate, Yes, DryRun bool
	KeepFrontmatter                                  bool
}

// DiscardFrontmatter reports whether frontmatter should be discarded
// during regeneration, based on the combination of --keep-frontmatter
// and the deprecated --force flag.
//
// Returns:
//   - bool: True if frontmatter should be discarded
func (o ExportOpts) DiscardFrontmatter() bool {
	return !o.KeepFrontmatter || o.Force
}

// FileAction describes the planned action for a single export file (one part
// of one session).
type FileAction struct {
	Session    *parser.Session
	Filename   string
	Path       string
	Part       int
	TotalParts int
	StartIdx   int
	EndIdx     int
	Action     ExportAction
	Messages   []parser.Message
	Slug       string
	Title      string
	BaseName   string
}

// ExportPlan is the result of PlanExport: a list of per-file actions plus
// aggregate counters and any renames that need to happen first.
type ExportPlan struct {
	Actions     []FileAction
	NewCount    int
	RegenCount  int
	SkipCount   int
	LockedCount int
	RenameOps   []RenameOp
}

// RenameOp describes a dedup rename (old slug → new slug).
type RenameOp struct {
	OldBase  string
	NewBase  string
	NumParts int
}

// PlanCounts converts an ExportPlan's counters to write.ExportCounts.
//
// Parameters:
//   - p: Export plan with counters
//
// Returns:
//   - write.ExportCounts: Formatted counters for output
func PlanCounts(p ExportPlan) write.ExportCounts {
	return write.ExportCounts{
		New:    p.NewCount,
		Regen:  p.RegenCount,
		Skip:   p.SkipCount,
		Locked: p.LockedCount,
	}
}
