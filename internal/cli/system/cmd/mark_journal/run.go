//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_journal

import (
	"github.com/spf13/cobra"

	coreJournal "github.com/ActiveMemory/ctx/internal/cli/system/core/journal"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	writeJournal "github.com/ActiveMemory/ctx/internal/write/mark_journal"
)

// Run handles the mark-journal command.
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
func Run(cmd *cobra.Command, filename, stage string) error {
	check, _ := cmd.Flags().GetBool(cFlag.Check)
	if check {
		r, checkErr := coreJournal.CheckStage(filename, stage)
		if checkErr != nil {
			return checkErr
		}
		writeJournal.StageChecked(cmd, filename, stage, r.Value)
		return nil
	}

	if markErr := coreJournal.MarkStage(filename, stage); markErr != nil {
		return markErr
	}

	writeJournal.StageMarked(cmd, filename, stage)
	return nil
}
