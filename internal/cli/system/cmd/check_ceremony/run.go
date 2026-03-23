//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_ceremony

import (
	"os"
	"path/filepath"

	ceremony2 "github.com/ActiveMemory/ctx/internal/cli/system/core/ceremony"
	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/ceremony"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the check-ceremonies hook logic.
//
// Scans recent journal files for /ctx-remember and /ctx-wrap-up usage. When
// either ceremony is missing, emits a nudge message and sends relay/nudge
// notifications. The check is daily-throttled and skipped when paused.
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

	input, _, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}

	remindedFile := filepath.Join(state.StateDir(), ceremony.CeremonyThrottleID)

	if hook2.DailyThrottled(remindedFile) {
		return nil
	}

	files := ceremony2.RecentJournalFiles(
		ctxContext.ResolvedJournalDir(), ceremony.CeremonyJournalLookback,
	)

	if len(files) == 0 {
		return nil
	}

	remember, wrapup := ceremony2.ScanJournalsForCeremonies(files)

	if remember && wrapup {
		return nil
	}

	msg, variant := ceremony2.EmitCeremonyNudge(remember, wrapup)
	writeHook.Nudge(cmd, msg)
	if msg == "" {
		return nil
	}
	ref := notify.NewTemplateRef(hook.CheckCeremonies, variant, nil)
	nudge.NudgeAndRelay(hook.CheckCeremonies+": "+
		desc.Text(text.DescKeyCeremonyRelayMessage),
		input.SessionID, ref,
	)
	internalIo.TouchFile(remindedFile)
	return nil
}
