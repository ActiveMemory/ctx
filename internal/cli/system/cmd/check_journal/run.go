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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	coreJournal "github.com/ActiveMemory/ctx/internal/cli/system/core/journal"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	ctxResolve "github.com/ActiveMemory/ctx/internal/context/resolve"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the check-journal hook logic.
//
// Checks for unimported Claude Code sessions and unenriched journal
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
	if !state.Initialized() {
		return nil
	}
	input, _, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	tmpDir := state.Dir()
	remindedFile := filepath.Join(tmpDir, journal.CheckJournalThrottleID)
	claudeProjectsDir := filepath.Join(
		os.Getenv(env.Home), journal.CheckJournalClaudeProjectsSubdir,
	)

	// Only remind once per day
	if coreCheck.DailyThrottled(remindedFile) {
		return nil
	}

	// Bail out if journal or Claude projects directories don't exist
	jDir := ctxResolve.ResolvedJournalDir()
	if _, statErr := os.Stat(jDir); os.IsNotExist(statErr) {
		return nil
	}
	if _, statErr := os.Stat(claudeProjectsDir); os.IsNotExist(statErr) {
		return nil
	}

	// Stage 1: Unimported sessions
	newestJournal := coreJournal.NewestMtime(jDir, file.ExtMarkdown)
	unimported := coreJournal.CountNewerFiles(
		claudeProjectsDir, file.ExtJSONL, newestJournal,
	)

	// Stage 2: Unenriched entries
	unenriched := coreJournal.CountUnenriched(jDir)

	if unimported == 0 && unenriched == 0 {
		return nil
	}

	vars := map[string]any{
		journal.VarUnimportedCount: unimported,
		journal.VarUnenrichedCount: unenriched,
	}

	var variant, fallback string
	switch {
	case unimported > 0 && unenriched > 0:
		variant = hook.VariantBoth
		fallback = fmt.Sprintf(desc.Text(
			text.DescKeyCheckJournalFallbackBoth), unimported, unenriched,
		)
	case unimported > 0:
		variant = hook.VariantUnimported
		fallback = fmt.Sprintf(desc.Text(
			text.DescKeyCheckJournalFallbackUnimported), unimported,
		)
	default:
		variant = hook.VariantUnenriched
		fallback = fmt.Sprintf(desc.Text(
			text.DescKeyCheckJournalFallbackUnenriched), unenriched,
		)
	}

	content := message.LoadMessage(hook.CheckJournal, variant, vars, fallback)
	if content == "" {
		return nil
	}

	boxTitle := desc.Text(text.DescKeyCheckJournalBoxTitle)
	relayPrefix := desc.Text(text.DescKeyCheckJournalRelayPrefix)

	writeHook.Nudge(cmd, message.NudgeBox(relayPrefix, boxTitle, content))

	ref := notify.NewTemplateRef(hook.CheckJournal, variant, vars)
	journalMsg := hook.CheckJournal + ": " + fmt.Sprintf(
		desc.Text(text.DescKeyCheckJournalRelayFormat),
		unimported, unenriched,
	)
	nudge.NudgeAndRelay(journalMsg, input.SessionID, ref)

	internalIo.TouchFile(remindedFile)
	return nil
}
