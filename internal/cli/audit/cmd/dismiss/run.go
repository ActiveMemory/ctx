//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/audit/core/store"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeAudit "github.com/ActiveMemory/ctx/internal/write/audit"
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
//
// Returns:
//   - error: non-nil on missing id, read, or ledger write
func Run(cmd *cobra.Command, ids []string, all bool) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}

	if all {
		count, dismissErr := store.DismissAll()
		if dismissErr != nil {
			return dismissErr
		}
		writeAudit.DismissedAll(cmd, count)
		return nil
	}

	for _, id := range ids {
		if dismissErr := store.Dismiss(id); dismissErr != nil {
			return dismissErr
		}
		writeAudit.Dismissed(cmd, id)
	}
	return nil
}
