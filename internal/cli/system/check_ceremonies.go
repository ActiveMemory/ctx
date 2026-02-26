//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
)

// checkCeremoniesCmd returns the "ctx system check-ceremonies" command.
//
// Scans recent journal entries for /ctx-remember and /ctx-wrap-up usage.
// If either is missing from the last 3 sessions, emits a VERBATIM relay
// nudge once per day encouraging adoption.
func checkCeremoniesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-ceremonies",
		Short: "Session ceremony nudge hook",
		Long: `Scans the last 3 journal entries for /ctx-remember and /ctx-wrap-up
usage. If either is missing, emits a VERBATIM relay nudge encouraging
adoption. Throttled to once per day.

Hook event: UserPromptSubmit
Output: VERBATIM relay (when ceremonies missing), silent otherwise
Silent when: both ceremonies found in recent sessions`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckCeremonies(cmd, os.Stdin)
		},
	}
}

func runCheckCeremonies(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}

	input := readInput(stdin)

	tmpDir := secureTempDir()
	remindedFile := filepath.Join(tmpDir, "ceremony-reminded")

	if isDailyThrottled(remindedFile) {
		return nil
	}

	files := recentJournalFiles(resolvedJournalDir(), 3)

	if len(files) == 0 {
		// No journal entries — skip ceremony check entirely.
		// The check-journal hook already nudges about missing exports.
		return nil
	}

	remember, wrapup := scanJournalsForCeremonies(files)

	if remember && wrapup {
		return nil
	}

	emitCeremonyNudge(cmd, remember, wrapup)
	_ = notify.Send("nudge", "check-ceremonies: Session ceremony nudge", input.SessionID, "")
	_ = notify.Send("relay", "check-ceremonies: Session ceremony nudge", input.SessionID, "")
	touchFile(remindedFile)
	return nil
}

// recentJournalFiles returns the n most recent .md files in the journal
// directory, sorted by filename descending (date prefix gives chronological
// order). Returns nil if the directory doesn't exist or has no .md files.
func recentJournalFiles(dir string, n int) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		names = append(names, e.Name())
	}

	// Sort descending — newest first (date-prefix filenames sort naturally).
	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	if len(names) > n {
		names = names[:n]
	}

	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = filepath.Join(dir, name)
	}
	return paths
}

// scanJournalsForCeremonies checks whether the given journal files contain
// references to /ctx-remember and /ctx-wrap-up.
func scanJournalsForCeremonies(files []string) (remember, wrapup bool) {
	for _, path := range files {
		data, err := os.ReadFile(path) //nolint:gosec // journal file path
		if err != nil {
			continue
		}
		content := string(data)
		if !remember && strings.Contains(content, "ctx-remember") {
			remember = true
		}
		if !wrapup && strings.Contains(content, "ctx-wrap-up") {
			wrapup = true
		}
		if remember && wrapup {
			return
		}
	}
	return
}

func emitCeremonyNudge(cmd *cobra.Command, remember, wrapup bool) {
	cmd.Println("IMPORTANT: Relay this session tip to the user VERBATIM before answering their question.")
	cmd.Println()

	dirFooter := contextDirLine()

	switch {
	case !remember && !wrapup:
		cmd.Println("┌─ Session Ceremonies ──────────────────────────────")
		cmd.Println("│ Your last 3 sessions didn't use /ctx-remember or")
		cmd.Println("│ /ctx-wrap-up.")
		cmd.Println("│")
		cmd.Println("│ Start sessions with:  /ctx-remember")
		cmd.Println("│   → Loads context, shows active tasks, picks up")
		cmd.Println("│     where you left off. No re-explaining needed.")
		cmd.Println("│")
		cmd.Println("│ End sessions with:    /ctx-wrap-up")
		cmd.Println("│   → Captures learnings and decisions so the next")
		cmd.Println("│     session starts informed, not from scratch.")
		cmd.Println("│")
		cmd.Println("│ These take seconds and save minutes.")
		if dirFooter != "" {
			cmd.Println("│ " + dirFooter)
		}
		cmd.Println("└───────────────────────────────────────────────────")

	case !remember:
		cmd.Println("┌─ Session Start ───────────────────────────────────")
		cmd.Println("│ Try starting this session with /ctx-remember")
		cmd.Println("│")
		cmd.Println("│ It loads your context, shows active tasks, and")
		cmd.Println("│ picks up where you left off — no re-explaining.")
		if dirFooter != "" {
			cmd.Println("│ " + dirFooter)
		}
		cmd.Println("└───────────────────────────────────────────────────")

	case !wrapup:
		cmd.Println("┌─ Session End ─────────────────────────────────────")
		cmd.Println("│ Your last 3 sessions didn't end with /ctx-wrap-up")
		cmd.Println("│")
		cmd.Println("│ It captures learnings and decisions so the next")
		cmd.Println("│ session starts informed, not from scratch.")
		if dirFooter != "" {
			cmd.Println("│ " + dirFooter)
		}
		cmd.Println("└───────────────────────────────────────────────────")
	}
}
