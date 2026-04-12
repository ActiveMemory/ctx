//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_ceremony

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCeremony "github.com/ActiveMemory/ctx/internal/cli/system/core/ceremony"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/ceremony"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	ctxResolve "github.com/ActiveMemory/ctx/internal/context/resolve"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
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

	input, _, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	remindedFile := filepath.Join(state.Dir(), ceremony.ThrottleID)

	if coreCheck.DailyThrottled(remindedFile) {
		return nil
	}

	files := coreCeremony.RecentJournalFiles(
		ctxResolve.JournalDir(), ceremony.JournalLookback,
	)

	if len(files) == 0 {
		return nil
	}

	remember, wrapUp := coreCeremony.ScanJournalsForCeremonies(files)

	if remember && wrapUp {
		return nil
	}

	msg, variant := coreCeremony.Emit(remember, wrapUp)
	writeSetup.Nudge(cmd, msg)
	if msg == "" {
		return nil
	}
	ref := notify.NewTemplateRef(hook.CheckCeremony, variant, nil)
	nudge.EmitAndRelay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckCeremony, desc.Text(text.DescKeyCeremonyRelayMessage)),
		input.SessionID, ref,
	)
	internalIo.TouchFile(remindedFile)
	return nil
}
