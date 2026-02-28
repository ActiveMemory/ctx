//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// checkJournalCmd returns the "ctx system check-journal" command.
//
// Detects unexported sessions and unenriched journal entries, then prints
// actionable commands. Runs once per day (throttled by marker file).
func checkJournalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-journal",
		Short: "Journal export/enrich reminder hook",
		Long: `Detects unexported Claude Code sessions and unenriched journal entries,
then prints actionable commands. Throttled to once per day.

Hook event: UserPromptSubmit
Output: VERBATIM relay with export/enrich commands, silent otherwise
Silent when: no unexported sessions and no unenriched entries`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckJournal(cmd, os.Stdin)
		},
	}
}

func runCheckJournal(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}
	input := readInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = sessionUnknown
	}
	if paused(sessionID) > 0 {
		return nil
	}

	tmpDir := secureTempDir()
	remindedFile := filepath.Join(tmpDir, "journal-reminded")
	claudeProjectsDir := filepath.Join(os.Getenv("HOME"), ".claude", "projects")

	// Only remind once per day
	if isDailyThrottled(remindedFile) {
		return nil
	}

	// Bail out if journal or Claude projects directories don't exist
	jDir := resolvedJournalDir()
	if _, err := os.Stat(jDir); os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(claudeProjectsDir); os.IsNotExist(err) {
		return nil
	}

	// Stage 1: Unexported sessions
	newestJournal := newestMtime(jDir, ".md")
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
		variant = variantBoth
		fallback = fmt.Sprintf("You have %d new session(s) not yet exported.\n", unexported) +
			fmt.Sprintf("%d existing entries need enrichment.\n", unenriched) +
			"\nExport and enrich:\n  ctx recall export --all\n  /ctx-journal-enrich-all"
	case unexported > 0:
		variant = "unexported"
		fallback = fmt.Sprintf("You have %d new session(s) not yet exported.\n", unexported) +
			"\nExport:\n  ctx recall export --all"
	default:
		variant = "unenriched"
		fallback = fmt.Sprintf("%d journal entries need enrichment.\n", unenriched) +
			"\nEnrich:\n  /ctx-journal-enrich-all"
	}

	content := loadMessage("check-journal", variant, vars, fallback)
	if content == "" {
		return nil
	}

	msg := "IMPORTANT: Relay this journal reminder to the user VERBATIM before answering their question.\n\n" +
		"┌─ Journal Reminder ─────────────────────────────\n"
	msg += boxLines(content)
	if line := contextDirLine(); line != "" {
		msg += "│ " + line + "\n"
	}
	msg += "└────────────────────────────────────────────────"
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-journal", variant, vars)
	journalMsg := fmt.Sprintf("check-journal: %d unexported, %d unenriched", unexported, unenriched)
	_ = notify.Send("nudge", journalMsg, input.SessionID, ref)
	_ = notify.Send("relay", journalMsg, input.SessionID, ref)
	eventlog.Append("relay", journalMsg, input.SessionID, ref)

	touchFile(remindedFile)
	return nil
}

// newestMtime returns the most recent mtime (as Unix timestamp) of files
// with the given extension in the directory. Returns 0 if none found.
func newestMtime(dir, ext string) int64 {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	var latest int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ext) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
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
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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
	jstate, err := state.Load(dir)
	if err != nil {
		return 0
	}
	return jstate.CountUnenriched(dir)
}
