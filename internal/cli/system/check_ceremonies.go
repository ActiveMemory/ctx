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

	msg := emitCeremonyNudge(cmd, remember, wrapup)
	if msg == "" {
		return nil
	}
	_ = notify.Send("nudge", "check-ceremonies: Session ceremony nudge", input.SessionID, msg)
	_ = notify.Send("relay", "check-ceremonies: Session ceremony nudge", input.SessionID, msg)
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

func emitCeremonyNudge(cmd *cobra.Command, remember, wrapup bool) string {
	var variant, boxTitle, fallback string

	switch {
	case !remember && !wrapup:
		variant = "both"
		boxTitle = "Session Ceremonies"
		fallback = "Your last 3 sessions didn't use /ctx-remember or\n" +
			"/ctx-wrap-up.\n" +
			"\n" +
			"Start sessions with:  /ctx-remember\n" +
			"  → Loads context, shows active tasks, picks up\n" +
			"    where you left off. No re-explaining needed.\n" +
			"\n" +
			"End sessions with:    /ctx-wrap-up\n" +
			"  → Captures learnings and decisions so the next\n" +
			"    session starts informed, not from scratch.\n" +
			"\n" +
			"These take seconds and save minutes."
	case !remember:
		variant = "remember"
		boxTitle = "Session Start"
		fallback = "Try starting this session with /ctx-remember\n" +
			"\n" +
			"It loads your context, shows active tasks, and\n" +
			"picks up where you left off — no re-explaining."
	case !wrapup:
		variant = "wrapup"
		boxTitle = "Session End"
		fallback = "Your last 3 sessions didn't end with /ctx-wrap-up\n" +
			"\n" +
			"It captures learnings and decisions so the next\n" +
			"session starts informed, not from scratch."
	}

	content := loadMessage("check-ceremonies", variant, nil, fallback)
	if content == "" {
		return ""
	}

	msg := "IMPORTANT: Relay this session tip to the user VERBATIM before answering their question.\n\n" +
		"┌─ " + boxTitle + " " + strings.Repeat("─", 51-len(boxTitle)) + "\n"
	msg += boxLines(content)
	if dirFooter := contextDirLine(); dirFooter != "" {
		msg += "│ " + dirFooter + "\n"
	}
	msg += "└───────────────────────────────────────────────────"

	cmd.Println(msg)
	return msg
}
