//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runSync scans all journal markdowns and syncs frontmatter lock state
// to .state.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on I/O failure
func runSync(cmd *cobra.Command) error {
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return fmt.Errorf("load journal state: %w", loadErr)
	}

	files, matchErr := core.MatchJournalFiles(journalDir, nil, true)
	if matchErr != nil {
		return matchErr
	}
	if len(files) == 0 {
		cmd.Println("No journal entries found.")
		return nil
	}

	locked, unlocked := 0, 0

	for _, filename := range files {
		path := filepath.Join(journalDir, filename)
		fmLocked := core.FrontmatterHasLocked(path)
		stateLocked := jstate.Locked(filename)

		switch {
		case fmLocked && !stateLocked:
			jstate.Mark(filename, "locked")
			cmd.Println(fmt.Sprintf("  ✓ %s (locked)", filename))
			locked++
		case !fmLocked && stateLocked:
			jstate.Clear(filename, "locked")
			cmd.Println(fmt.Sprintf("  ✓ %s (unlocked)", filename))
			unlocked++
		}
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return fmt.Errorf("save journal state: %w", saveErr)
	}

	if locked == 0 && unlocked == 0 {
		cmd.Println("No changes — state already matches frontmatter.")
	} else {
		if locked > 0 {
			cmd.Println(fmt.Sprintf("\nLocked %d entry(s).", locked))
		}
		if unlocked > 0 {
			cmd.Println(fmt.Sprintf("\nUnlocked %d entry(s).", unlocked))
		}
	}

	return nil
}
