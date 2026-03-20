//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_journal

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	ctxcontext "github.com/ActiveMemory/ctx/internal/context/resolve"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/journal/state"
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
	journalDir := ctxcontext.ResolvedJournalDir()

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return ctxerr.LoadStateFailed(loadErr)
	}

	check, _ := cmd.Flags().GetBool("check")
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
			return ctxerr.UnknownStage(stage, strings.Join(state.ValidStages, ", "))
		}
		if val == "" {
			return ctxerr.StageNotSet(filename, stage)
		}
		cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyMarkJournalChecked), filename, stage, val))
		return nil
	}

	if ok := jstate.Mark(filename, stage); !ok {
		return ctxerr.UnknownStage(stage, strings.Join(state.ValidStages, ", "))
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return ctxerr.SaveStateFailed(saveErr)
	}

	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyMarkJournalMarked), filename, stage))
	return nil
}
