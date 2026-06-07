//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
)

// Dream execution mode constants.
const (
	// ModeDiscipline triages ideas/ by the discipline clock (hash
	// change). The only mode built in v1.
	ModeDiscipline Mode = "discipline"
	// ModeCreative is the deferred resurfacing mode (sketched, not
	// built in v1).
	ModeCreative Mode = "creative"
)

// Numeric defaults for a dream pass when neither flag nor rc
// supplies a value.
const (
	// DefaultMax is the default ceiling on ideas/ files processed
	// per pass.
	DefaultMax = 50
	// DefaultBudget is the default step/token budget for a pass.
	DefaultBudget = 30
	// DefaultQuietMinutes is the default activity quiet window the
	// trigger gate honors.
	DefaultQuietMinutes = 60
)

// Notebook artifacts and layout names within the gitignored
// dreams/ directory.
const (
	// FileLock is the flock lock file under dreams/ that serializes
	// passes.
	FileLock = ".lock"
	// FileFailed is the failmark a fail-loud pass leaves under
	// dreams/ when the executor cannot run.
	FileFailed = ".failed"
	// FileProposals is the proposals file the executor writes into a
	// per-run dreams/<ts>/ directory.
	FileProposals = "proposals.json"
	// BackupSuffix is appended to a backed-up source file's base
	// name before a destructive mutation.
	BackupSuffix = ".bak"
	// RunTimeLayout is the timestamp layout for a per-run
	// dreams/<ts>/ directory name (UTC, compact). It reuses the
	// canonical compact RFC3339 layout.
	RunTimeLayout = cfgTime.RFC3339Compact
)

// IdeaGlob is the markdown extension the scanner matches under
// ideas/ (its own dreams/ notebook and binaries are excluded).
const IdeaGlob = cfgFile.ExtMarkdown

// BlogMarker is the line appended in place to tag an idea as blog
// material for the mark-blog mechanical disposition.
const BlogMarker = "<!-- ctx-dream: blog-candidate -->"

// Executor defaults. The reference executor is a headless
// claude -p invocation; an empty rc value selects this default.
const (
	// ExecutorDefaultBin is the reference executor binary.
	ExecutorDefaultBin = "claude"
	// ExecutorPromptFlag is the headless prompt flag for the
	// reference executor.
	ExecutorPromptFlag = "-p"
	// ExecutorMaxTurnsFlag bounds the reference executor's agentic
	// loop to the pass budget.
	ExecutorMaxTurnsFlag = "--max-turns"
	// ExecutorModelFlag selects the executor model when a model
	// override is configured (empty uses the executor's default).
	ExecutorModelFlag = "--model"
	// ExecutorPromptTemplate instructs the executor to run the
	// ctx-dream skill with the pass parameters and notebook paths.
	// Order: mode, max, budget, ideas dir, dreams run dir.
	ExecutorPromptTemplate = "Run the ctx-dream skill in %s mode over " +
		"at most %d ideas files with a step budget of %d. Read ideas " +
		"from %s and write provenance-bearing proposals only into %s. " +
		"Never touch canonical memory."
)
