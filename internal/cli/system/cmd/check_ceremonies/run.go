//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_ceremonies

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
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
	if !core.IsInitialized() {
		return nil
	}

	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	remindedFile := filepath.Join(core.StateDir(), config.CeremonyThrottleID)

	if core.IsDailyThrottled(remindedFile) {
		return nil
	}

	files := core.RecentJournalFiles(
		core.ResolvedJournalDir(), config.CeremonyJournalLookback,
	)

	if len(files) == 0 {
		return nil
	}

	remember, wrapup := core.ScanJournalsForCeremonies(files)

	if remember && wrapup {
		return nil
	}

	msg, variant := core.EmitCeremonyNudge(cmd, remember, wrapup)
	if msg == "" {
		return nil
	}
	ref := notify.NewTemplateRef(config.HookCheckCeremonies, variant, nil)
	core.NudgeAndRelay(config.HookCheckCeremonies+": "+
		assets.TextDesc(assets.TextDescKeyCeremonyRelayMessage),
		input.SessionID, ref,
	)
	core.TouchFile(remindedFile)
	return nil
}
