//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_journal

import (
	"strings"

	"github.com/spf13/cobra"

	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxResolve "github.com/ActiveMemory/ctx/internal/context/resolve"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	writeJournal "github.com/ActiveMemory/ctx/internal/write/mark_journal"
)

// runMarkJournal handles the mark-journal command.
//
// Marks a journal file as having reached a given processing stage, or
// checks the current stage value when --check is set.
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//   - filename: journal filename to mark or check
//   - stage: processing stage name (exported, enriched, normalized, etc.)
//
// Returns:
//   - error: Non-nil on state load/save failure or unknown stage
func runMarkJournal(cmd *cobra.Command, filename, stage string) error {
	journalDir := ctxResolve.ResolvedJournalDir()

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return errJournal.LoadStateFailed(loadErr)
	}

	check, _ := cmd.Flags().GetBool(cflag.Check)
	if check {
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
			return errJournal.UnknownStage(
				stage, strings.Join(state.ValidStages, token.CommaSpace),
			)
		}
		if val == "" {
			return errJournal.StageNotSet(filename, stage)
		}
		writeJournal.StageChecked(cmd, filename, stage, val)
		return nil
	}

	if ok := jstate.Mark(filename, stage); !ok {
		return errJournal.UnknownStage(
			stage, strings.Join(state.ValidStages, token.CommaSpace),
		)
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return errJournal.SaveStateFailed(saveErr)
	}

	writeJournal.StageMarked(cmd, filename, stage)
	return nil
}
