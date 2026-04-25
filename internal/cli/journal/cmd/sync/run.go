//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/lock"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeRecall "github.com/ActiveMemory/ctx/internal/write/journal"
)

// Run scans all journal markdowns and syncs frontmatter lock state
// to .state.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on I/O failure
func Run(cmd *cobra.Command) error {
	ctxDir, ctxErr := rc.RequireContextDir()
	if ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	journalDir := filepath.Join(ctxDir, dir.Journal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return errJournal.LoadState(loadErr)
	}

	files, matchErr := lock.MatchJournalFiles(journalDir, nil, true)
	if matchErr != nil {
		return matchErr
	}
	if len(files) == 0 {
		writeRecall.SyncNone(cmd)
		return nil
	}

	locked, unlocked := 0, 0

	for _, filename := range files {
		path := filepath.Join(journalDir, filename)
		fmLocked := lock.HasLocked(path)
		stateLocked := jstate.Locked(filename)

		switch {
		case fmLocked && !stateLocked:
			jstate.Mark(filename, journal.StageLocked)
			writeRecall.SyncLocked(cmd, filename)
			locked++
		case !fmLocked && stateLocked:
			jstate.Clear(filename, journal.StageLocked)
			writeRecall.SyncUnlocked(cmd, filename)
			unlocked++
		}
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return errJournal.SaveState(saveErr)
	}

	writeRecall.SyncSummary(cmd, locked, unlocked)

	return nil
}
