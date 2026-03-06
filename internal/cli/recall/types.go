//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
)

// exportAction describes what will happen to a given file.
type exportAction int

const (
	actionNew        exportAction = iota // file does not exist yet
	actionRegenerate                     // file exists and will be rewritten
	actionSkip                           // file exists and will be left alone
	actionLocked                         // file is locked — never overwritten
)

// exportOpts holds all flag values for the export command.
type exportOpts struct {
	all, allProjects, force, regenerate, yes, dryRun bool
	keepFrontmatter                                  bool
}

// discardFrontmatter reports whether frontmatter should be discarded
// during regeneration, based on the combination of --keep-frontmatter
// and the deprecated --force flag.
func (o exportOpts) discardFrontmatter() bool {
	return !o.keepFrontmatter || o.force
}

// fileAction describes the planned action for a single export file (one part
// of one session).
type fileAction struct {
	session    *parser.Session
	filename   string
	path       string
	part       int
	totalParts int
	startIdx   int
	endIdx     int
	action     exportAction
	messages   []parser.Message
	slug       string
	title      string
	baseName   string
}

// exportPlan is the result of planExport: a list of per-file actions plus
// aggregate counters and any renames that need to happen first.
type exportPlan struct {
	actions     []fileAction
	newCount    int
	regenCount  int
	skipCount   int
	lockedCount int
	renameOps   []renameOp
}

// renameOp describes a dedup rename (old slug → new slug).
type renameOp struct {
	oldBase  string
	newBase  string
	numParts int
}

// planCounts converts an exportPlan's counters to write.ExportCounts.
func planCounts(p exportPlan) write.ExportCounts {
	return write.ExportCounts{
		New:    p.newCount,
		Regen:  p.regenCount,
		Skip:   p.skipCount,
		Locked: p.lockedCount,
	}
}
