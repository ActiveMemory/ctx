//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkaudit

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreProv "github.com/ActiveMemory/ctx/internal/cli/system/core/provenance"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	auditStore "github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/core/store"
	cfgAudit "github.com/ActiveMemory/ctx/internal/ctxctl/config/audit"
	auditRender "github.com/ActiveMemory/ctx/internal/ctxctl/core/audit"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Run executes the audit-relay hook logic.
//
// Reads hook input from stdin, loads every audit report,
// filters by status (skip clean) and dismissal (skip
// dismissed-against-current-digest), then emits a single
// verbatim-relay box concatenating each remaining
// report's body. Stale reports (older than
// [cfgAudit.StalenessAge]) are prefixed with a STALE
// marker but are still relayed.
//
// The relay box copy is supplied verbatim by ctxctl via s;
// this logic holds no user-facing text of its own and does
// not route through ctx's hook-message template machinery.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//   - s: English user-facing text supplied by ctxctl
//
// Returns:
//   - error: nil on intentional silence; propagated from the
//     nudge emit on a real relay-log or webhook failure
func Run(cmd *cobra.Command, stdin *os.File, s Strings) error {
	input, _, paused := coreCheck.Preamble(stdin)

	coreProv.Emit(cmd, input.SessionID)

	initialized, initErr := state.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return nil
	}
	if !initialized || paused {
		return nil
	}

	reports, readErr := auditStore.Read()
	if readErr != nil || len(reports) == 0 {
		return nil
	}
	led, ledErr := auditStore.ReadDismissals()
	if ledErr != nil {
		return nil
	}

	var due []auditStore.Report
	for _, r := range reports {
		if r.Status == cfgAudit.StatusClean {
			continue
		}
		if auditStore.IsDismissed(r, led) {
			continue
		}
		due = append(due, r)
	}
	if len(due) == 0 {
		return nil
	}

	body := auditRender.RenderReports(
		due, time.Now().UTC(), s.ReportSeparator, s.StalePrefix,
	)

	content := body +
		token.NewlineLF + s.DismissHint +
		token.NewlineLF + s.DismissAllHint
	relayMsg := fmt.Sprintf(s.NudgeFormat, len(due))

	return nudge.Emit(cmd, content,
		s.RelayPrefix, s.BoxTitle,
		s.RelayLabel, s.RelayVariant,
		relayMsg, input.SessionID, nil, "",
	)
}
