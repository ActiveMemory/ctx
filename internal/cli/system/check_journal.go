//   /    Context:                     https://ctx.ist
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

	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/spf13/cobra"
)

// checkJournalCmd returns the "ctx system check-journal" command.
//
// Detects unexported sessions and unenriched journal entries, then prints
// actionable commands. Runs once per day (throttled by marker file).
func checkJournalCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "check-journal",
		Short:  "Journal export/enrich reminder hook",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckJournal(cmd)
		},
	}
}

func runCheckJournal(cmd *cobra.Command) error {
	if !isInitialized() {
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
	if _, err := os.Stat(journalDir); os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(claudeProjectsDir); os.IsNotExist(err) {
		return nil
	}

	// Stage 1: Unexported sessions
	newestJournal := newestMtime(journalDir, ".md")
	unexported := countNewerFiles(claudeProjectsDir, ".jsonl", newestJournal)

	// Stage 2: Unenriched entries
	unenriched := countUnenriched(journalDir)

	if unexported == 0 && unenriched == 0 {
		return nil
	}

	cmd.Println("IMPORTANT: Relay this journal reminder to the user VERBATIM before answering their question.")
	cmd.Println()
	cmd.Println("┌─ Journal Reminder ─────────────────────────────")

	switch {
	case unexported > 0 && unenriched > 0:
		cmd.Println(fmt.Sprintf("│ You have %d new session(s) not yet exported.", unexported))
		cmd.Println(fmt.Sprintf("│ %d existing entries need enrichment.", unenriched))
		cmd.Println("│")
		cmd.Println("│ Export, normalize, then enrich:")
		cmd.Println("│   ctx recall export --all")
		cmd.Println("│   /ctx-journal-normalize")
		cmd.Println("│   /ctx-journal-enrich-all")
	case unexported > 0:
		cmd.Println(fmt.Sprintf("│ You have %d new session(s) not yet exported.", unexported))
		cmd.Println("│")
		cmd.Println("│ Export:")
		cmd.Println("│   ctx recall export --all")
	default:
		cmd.Println(fmt.Sprintf("│ %d journal entries need enrichment.", unenriched))
		cmd.Println("│")
		cmd.Println("│ Normalize, then enrich:")
		cmd.Println("│   /ctx-journal-normalize")
		cmd.Println("│   /ctx-journal-enrich-all")
	}

	cmd.Println("└────────────────────────────────────────────────")

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
