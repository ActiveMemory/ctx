//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_journal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	ctxcontext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Run executes the check-journal hook logic.
//
// Checks for unexported Claude Code sessions and unenriched journal
// entries, then emits a journal reminder nudge if either is found.
// Throttled to once per day.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.Initialized() {
		return nil
	}
	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	tmpDir := core.StateDir()
	remindedFile := filepath.Join(tmpDir, journal.CheckJournalThrottleID)
	claudeProjectsDir := filepath.Join(
		os.Getenv(env.Home), journal.CheckJournalClaudeProjectsSubdir,
	)

	// Only remind once per day
	if core.IsDailyThrottled(remindedFile) {
		return nil
	}

	// Bail out if journal or Claude projects directories don't exist
	jDir := ctxcontext.ResolvedJournalDir()
	if _, statErr := os.Stat(jDir); os.IsNotExist(statErr) {
		return nil
	}
	if _, statErr := os.Stat(claudeProjectsDir); os.IsNotExist(statErr) {
		return nil
	}

	// Stage 1: Unexported sessions
	newestJournal := core.NewestMtime(jDir, file.ExtMarkdown)
	unexported := core.CountNewerFiles(
		claudeProjectsDir, file.ExtJSONL, newestJournal,
	)

	// Stage 2: Unenriched entries
	unenriched := core.CountUnenriched(jDir)

	if unexported == 0 && unenriched == 0 {
		return nil
	}

	vars := map[string]any{
		tpl.VarUnexportedCount: unexported,
		tpl.VarUnenrichedCount: unenriched,
	}

	var variant, fallback string
	switch {
	case unexported > 0 && unenriched > 0:
		variant = hook.VariantBoth
		fallback = fmt.Sprintf(assets.TextDesc(
			assets.TextDescKeyCheckJournalFallbackBoth), unexported, unenriched,
		)
	case unexported > 0:
		variant = hook.VariantUnexported
		fallback = fmt.Sprintf(assets.TextDesc(
			assets.TextDescKeyCheckJournalFallbackUnexported), unexported,
		)
	default:
		variant = hook.VariantUnenriched
		fallback = fmt.Sprintf(assets.TextDesc(
			assets.TextDescKeyCheckJournalFallbackUnenriched), unenriched,
		)
	}

	content := core.LoadMessage(hook.CheckJournal, variant, vars, fallback)
	if content == "" {
		return nil
	}

	boxTitle := assets.TextDesc(assets.TextDescKeyCheckJournalBoxTitle)
	relayPrefix := assets.TextDesc(assets.TextDescKeyCheckJournalRelayPrefix)

	cmd.Println(core.NudgeBox(relayPrefix, boxTitle, content))

	ref := notify.NewTemplateRef(hook.CheckJournal, variant, vars)
	journalMsg := hook.CheckJournal + ": " + fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyCheckJournalRelayFormat),
		unexported, unenriched,
	)
	core.NudgeAndRelay(journalMsg, input.SessionID, ref)

	core.TouchFile(remindedFile)
	return nil
}
