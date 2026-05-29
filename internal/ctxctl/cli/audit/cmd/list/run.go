//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"

	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/core/store"
	cfgAudit "github.com/ActiveMemory/ctx/internal/ctxctl/config/audit"
	writeAudit "github.com/ActiveMemory/ctx/internal/ctxctl/write/audit"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run lists every audit report under `.context/audit/`
// with status, commit-range, and generated-at timestamp.
// Dismissal state is annotated by suffixing the status
// with "(dismissed)" so a single column carries the
// hook-relevant filter signal.
//
// Parameters:
//   - cmd: Cobra command for output
//   - noneMsg: message printed when no reports exist
//   - itemFormat: per-row format string (id, status,
//     commit-range, generated-at)
//
// Returns:
//   - error: non-nil on read or ledger I/O failure
func Run(cmd *cobra.Command, noneMsg, itemFormat string) error {
	// Every error from here is actionable, not a usage problem, so
	// suppress cobra's help dump (the ctxctl root silences cobra's
	// error line; printErr is the sole printer).
	cmd.SilenceUsage = true
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		return ctxErr
	}

	reports, readErr := store.Read()
	if readErr != nil {
		return readErr
	}
	if len(reports) == 0 {
		writeAudit.None(cmd, noneMsg)
		return nil
	}

	led, ledErr := store.ReadDismissals()
	if ledErr != nil {
		return ledErr
	}

	for _, r := range reports {
		status := r.Status
		if store.IsDismissed(r, led) {
			status += cfgAudit.SuffixDismissed
		}
		writeAudit.ListItem(cmd,
			itemFormat, r.ID, status, r.CommitRange,
			r.GeneratedAt.Format(cfgTime.DateTimeFmt),
		)
	}
	return nil
}
