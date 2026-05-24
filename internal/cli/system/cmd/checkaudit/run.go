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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	auditStore "github.com/ActiveMemory/ctx/internal/cli/audit/core/store"
	auditRender "github.com/ActiveMemory/ctx/internal/cli/system/core/audit"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreProv "github.com/ActiveMemory/ctx/internal/cli/system/core/provenance"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	cfgAudit "github.com/ActiveMemory/ctx/internal/config/audit"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Run executes the check-audit hook logic.
//
// Reads hook input from stdin, loads every audit report,
// filters by status (skip clean) and dismissal (skip
// dismissed-against-current-digest), then emits a single
// verbatim-relay box concatenating each remaining
// report's body. Stale reports (older than
// [cfgAudit.StalenessAge]) are prefixed with a STALE
// marker but are still relayed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
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

	body := auditRender.RenderReports(due, time.Now().UTC())

	fallback := body +
		token.NewlineLF +
		desc.Text(text.DescKeyCheckAuditDismissHint) +
		token.NewlineLF +
		desc.Text(text.DescKeyCheckAuditDismissAllHint)
	vars := map[string]any{cfgAudit.VarList: body}
	relayMsg := fmt.Sprintf(
		desc.Text(text.DescKeyCheckAuditNudgeFormat),
		len(due),
	)
	return nudge.LoadAndEmit(cmd,
		hook.CheckAudit, hook.VariantAudits,
		vars, fallback,
		desc.Text(text.DescKeyCheckAuditRelayPrefix),
		desc.Text(text.DescKeyCheckAuditBoxTitle),
		relayMsg, input.SessionID, "",
	)
}
