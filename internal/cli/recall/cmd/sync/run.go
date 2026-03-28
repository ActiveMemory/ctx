//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/lock"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/recall"
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
	journalDir := filepath.Join(rc.ContextDir(), dir.Journal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return errJournal.LoadState(loadErr)
	}

	files, matchErr := lock.MatchJournalFiles(journalDir, nil, true)
	if matchErr != nil {
		return matchErr
	}
	if len(files) == 0 {
		recall.JournalSyncNone(cmd)
		return nil
	}

	locked, unlocked := 0, 0

	for _, filename := range files {
		path := filepath.Join(journalDir, filename)
		fmLocked := lock.FrontmatterHasLocked(path)
		stateLocked := jstate.Locked(filename)

		switch {
		case fmLocked && !stateLocked:
			jstate.Mark(filename, journal.StageLocked)
			recall.JournalSyncLocked(cmd, filename)
			locked++
		case !fmLocked && stateLocked:
			jstate.Clear(filename, journal.StageLocked)
			recall.JournalSyncUnlocked(cmd, filename)
			unlocked++
		}
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return errJournal.SaveState(saveErr)
	}

	recall.JournalSyncSummary(cmd, locked, unlocked)

	return nil
}
