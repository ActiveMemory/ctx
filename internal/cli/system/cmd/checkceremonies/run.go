//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkceremonies

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

func runCheckCeremonies(cmd *cobra.Command, stdin *os.File) error {
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
	remindedFile := filepath.Join(tmpDir, "ceremony-reminded")

	if core.IsDailyThrottled(remindedFile) {
		return nil
	}

	files := recentJournalFiles(core.ResolvedJournalDir(), 3)

	if len(files) == 0 {
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
	var variant string
	switch {
	case !remember && !wrapup:
		variant = core.VariantBoth
	case !remember:
		variant = "remember"
	default:
		variant = "wrapup"
	}
	ref := notify.NewTemplateRef("check-ceremonies", variant, nil)
	_ = notify.Send("nudge", "check-ceremonies: Session ceremony nudge", input.SessionID, ref)
	_ = notify.Send("relay", "check-ceremonies: Session ceremony nudge", input.SessionID, ref)
	eventlog.Append("relay", "check-ceremonies: Session ceremony nudge", input.SessionID, ref)
	core.TouchFile(remindedFile)
	return nil
}

// recentJournalFiles returns the n most recent .md files in the journal
// directory, sorted by filename descending.
func recentJournalFiles(dir string, n int) []string {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), config.ExtMarkdown) {
			continue
		}
		names = append(names, e.Name())
	}

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
		data, readErr := os.ReadFile(path) //nolint:gosec // journal file path
		if readErr != nil {
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
		variant = core.VariantBoth
		boxTitle = "Session Ceremonies"
		fallback = "Your last 3 sessions didn't use /ctx-remember or\n" +
			"/ctx-wrap-up.\n" +
			config.NewlineLF +
			"Start sessions with:  /ctx-remember\n" +
			"  → Loads context, shows active tasks, picks up\n" +
			"    where you left off. No re-explaining needed.\n" +
			config.NewlineLF +
			"End sessions with:    /ctx-wrap-up\n" +
			"  → Captures learnings and decisions so the next\n" +
			"    session starts informed, not from scratch.\n" +
			config.NewlineLF +
			"These take seconds and save minutes."
	case !remember:
		variant = "remember"
		boxTitle = "Session Start"
		fallback = "Try starting this session with /ctx-remember\n" +
			config.NewlineLF +
			"It loads your context, shows active tasks, and\n" +
			"picks up where you left off — no re-explaining."
	case !wrapup:
		variant = "wrapup"
		boxTitle = "Session End"
		fallback = "Your last 3 sessions didn't end with /ctx-wrap-up\n" +
			config.NewlineLF +
			"It captures learnings and decisions so the next\n" +
			"session starts informed, not from scratch."
	}

	content := core.LoadMessage("check-ceremonies", variant, nil, fallback)
	if content == "" {
		return ""
	}

	msg := "IMPORTANT: Relay this session tip to the user VERBATIM before answering their question.\n\n" +
		"┌─ " + boxTitle + " " + strings.Repeat("─", 51-len(boxTitle)) + config.NewlineLF
	msg += core.BoxLines(content)
	if dirFooter := core.ContextDirLine(); dirFooter != "" {
		msg += "│ " + dirFooter + config.NewlineLF
	}
	msg += "└───────────────────────────────────────────────────"

	cmd.Println(msg)
	return msg
}
