//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/core/store"
	writeAudit "github.com/ActiveMemory/ctx/internal/ctxctl/write/audit"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run dismisses one or more audit reports. With all=true,
// every report currently on disk is dismissed; otherwise
// each id in ids is dismissed individually and a missing
// id surfaces as [errAudit.UnknownID].
//
// Parameters:
//   - cmd: Cobra command for output
//   - ids: report basenames (ignored when all is true)
//   - all: when true, dismiss every report
//   - dismissedFmt: single-dismissal confirmation format (id)
//   - dismissedAllFmt: bulk-dismissal confirmation format (count)
//
// Returns:
//   - error: non-nil on missing id, read, or ledger write
func Run(
	cmd *cobra.Command, ids []string, all bool,
	dismissedFmt, dismissedAllFmt string,
) error {
	// Every error from here is actionable, not a usage problem, so
	// suppress cobra's help dump (the ctxctl root silences cobra's
	// error line; printErr is the sole printer).
	cmd.SilenceUsage = true
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		return ctxErr
	}

	if all {
		count, dismissErr := store.DismissAll()
		if dismissErr != nil {
			return dismissErr
		}
		writeAudit.DismissedAll(cmd, dismissedAllFmt, count)
		return nil
	}

	for _, id := range ids {
		if dismissErr := store.Dismiss(id); dismissErr != nil {
			return dismissErr
		}
		writeAudit.Dismissed(cmd, dismissedFmt, id)
	}
	return nil
}
