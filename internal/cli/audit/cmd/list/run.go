//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/audit/core/store"
	cfgAudit "github.com/ActiveMemory/ctx/internal/config/audit"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeAudit "github.com/ActiveMemory/ctx/internal/write/audit"
)

// Run lists every audit report under `.context/audit/`
// with status, commit-range, and generated-at timestamp.
// Dismissal state is annotated by suffixing the status
// with "(dismissed)" so a single column carries the
// hook-relevant filter signal.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: non-nil on read or ledger I/O failure
func Run(cmd *cobra.Command) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}

	reports, readErr := store.Read()
	if readErr != nil {
		return readErr
	}
	if len(reports) == 0 {
		writeAudit.None(cmd)
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
			r.ID, status, r.CommitRange,
			r.GeneratedAt.Format(cfgTime.DateTimeFmt),
		)
	}
	return nil
}
