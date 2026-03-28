//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxResolve "github.com/ActiveMemory/ctx/internal/context/resolve"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// CheckResult holds the outcome of checking a stage value.
type CheckResult struct {
	Value string
}

// MarkResult holds the outcome of marking a stage.
type MarkResult struct {
	Marked bool
}

// CheckStage reads the current value of a processing stage for a file.
//
// Parameters:
//   - filename: journal filename to check
//   - stage: processing stage name
//
// Returns:
//   - CheckResult: Stage value
//   - error: Non-nil on state load failure, unknown stage, or unset stage
func CheckStage(filename, stage string) (CheckResult, error) {
	journalDir := ctxResolve.ResolvedJournalDir()
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return CheckResult{}, errJournal.LoadStateFailed(loadErr)
	}

	fs := jstate.Entries[filename]
	var val string
	switch stage {
	case journal.StageExported:
		val = fs.Exported
	case journal.StageEnriched:
		val = fs.Enriched
	case journal.StageNormalized:
		val = fs.Normalized
	case journal.StageFencesVerified:
		val = fs.FencesVerified
	case journal.StageLocked:
		val = fs.Locked
	default:
		return CheckResult{}, errJournal.UnknownStage(
			stage, strings.Join(state.ValidStages, token.CommaSpace),
		)
	}
	if val == "" {
		return CheckResult{}, errJournal.StageNotSet(filename, stage)
	}
	return CheckResult{Value: val}, nil
}

// MarkStage sets a processing stage for a journal file and persists state.
//
// Parameters:
//   - filename: journal filename to mark
//   - stage: processing stage name
//
// Returns:
//   - error: Non-nil on unknown stage or state load/save failure
func MarkStage(filename, stage string) error {
	journalDir := ctxResolve.ResolvedJournalDir()
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return errJournal.LoadStateFailed(loadErr)
	}

	if ok := jstate.Mark(filename, stage); !ok {
		return errJournal.UnknownStage(
			stage, strings.Join(state.ValidStages, token.CommaSpace),
		)
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return errJournal.SaveStateFailed(saveErr)
	}

	return nil
}
