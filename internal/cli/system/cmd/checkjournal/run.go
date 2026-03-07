//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkjournal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/notify"
)

func runCheckJournal(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}
	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}

	tmpDir := core.StateDir()
	remindedFile := filepath.Join(tmpDir, "journal-reminded")
	claudeProjectsDir := filepath.Join(os.Getenv("HOME"), ".claude", "projects")

	// Only remind once per day
	if core.IsDailyThrottled(remindedFile) {
		return nil
	}

	// Bail out if journal or Claude projects directories don't exist
	jDir := core.ResolvedJournalDir()
	if _, statErr := os.Stat(jDir); os.IsNotExist(statErr) {
		return nil
	}
	if _, statErr := os.Stat(claudeProjectsDir); os.IsNotExist(statErr) {
		return nil
	}

	// Stage 1: Unexported sessions
	newestJournal := newestMtime(jDir, config.ExtMarkdown)
	unexported := countNewerFiles(claudeProjectsDir, ".jsonl", newestJournal)

	// Stage 2: Unenriched entries
	unenriched := countUnenriched(jDir)

	if unexported == 0 && unenriched == 0 {
		return nil
	}

	vars := map[string]any{
		"UnexportedCount": unexported,
		"UnenrichedCount": unenriched,
	}

	var variant, fallback string
	switch {
	case unexported > 0 && unenriched > 0:
		variant = core.VariantBoth
		fallback = fmt.Sprintf("You have %d new session(s) not yet exported.\n", unexported) +
			fmt.Sprintf("%d existing entries need enrichment.\n", unenriched) +
			"\nProcess journal (exports and enriches):\n  /ctx-journal-enrich-all"
	case unexported > 0:
		variant = "unexported"
		fallback = fmt.Sprintf("You have %d new session(s) not yet exported.\n", unexported) +
			"\nProcess journal (exports and enriches):\n  /ctx-journal-enrich-all"
	default:
		variant = "unenriched"
		fallback = fmt.Sprintf("%d journal entries need enrichment.\n", unenriched) +
			"\nEnrich:\n  /ctx-journal-enrich-all"
	}

	content := core.LoadMessage("check-journal", variant, vars, fallback)
	if content == "" {
		return nil
	}

	msg := "IMPORTANT: Relay this journal reminder to the user VERBATIM before answering their question.\n\n" +
		"┌─ Journal Reminder ─────────────────────────────\n"
	msg += core.BoxLines(content)
	if line := core.ContextDirLine(); line != "" {
		msg += "│ " + line + config.NewlineLF
	}
	msg += "└────────────────────────────────────────────────"
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-journal", variant, vars)
	journalMsg := fmt.Sprintf("check-journal: %d unexported, %d unenriched", unexported, unenriched)
	_ = notify.Send("nudge", journalMsg, input.SessionID, ref)
	_ = notify.Send("relay", journalMsg, input.SessionID, ref)
	eventlog.Append("relay", journalMsg, input.SessionID, ref)

	core.TouchFile(remindedFile)
	return nil
}

// newestMtime returns the most recent mtime (as Unix timestamp) of files
// with the given extension in the directory. Returns 0 if none found.
func newestMtime(dir, ext string) int64 {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return 0
	}

	var latest int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ext) {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		mtime := info.ModTime().Unix()
		if mtime > latest {
			latest = mtime
		}
	}
	return latest
}

// countNewerFiles recursively counts files with the given extension that
// are newer than the reference timestamp.
func countNewerFiles(dir, ext string, refTime int64) int {
	count := 0
	_ = filepath.Walk(dir, func(_ string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil // skip errors
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ext) {
			return nil
		}
		if info.ModTime().Unix() > refTime {
			count++
		}
		return nil
	})
	return count
}

// countUnenriched counts journal .md files that lack an enriched date
// in the journal state file.
func countUnenriched(dir string) int {
	jstate, loadErr := state.Load(dir)
	if loadErr != nil {
		return 0
	}
	return jstate.CountUnenriched(dir)
}
