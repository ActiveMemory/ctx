//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package review

import (
	"github.com/spf13/cobra"

	dreamPaths "github.com/ActiveMemory/ctx/internal/cli/dream/core/paths"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	writeDream "github.com/ActiveMemory/ctx/internal/write/dream"
)

// Run lists the pending proposals from the latest dream run and
// renders them. Proposals already recorded in the ledger are filtered
// out (dedup-against-seen). An absent run or empty pending set prints
// the no-pending message.
//
// Parameters:
//   - cmd: cobra command for output
//
// Returns:
//   - error: a resolution or read failure
func Run(cmd *cobra.Command) error {
	loc, locErr := dreamPaths.Resolve()
	if locErr != nil {
		return locErr
	}
	runDir, runErr := engine.LatestRunDir(loc.Dreams)
	if runErr != nil {
		return runErr
	}
	if runDir == "" {
		writeDream.ReviewNone(cmd)
		return nil
	}
	proposals, readErr := engine.ReadProposals(runDir)
	if readErr != nil {
		return readErr
	}
	ledger, ledgerErr := engine.ReadLedger(loc.Dreams)
	if ledgerErr != nil {
		return ledgerErr
	}
	pending := engine.PendingProposals(proposals, ledger)
	if len(pending) == 0 {
		writeDream.ReviewNone(cmd)
		return nil
	}
	writeDream.Review(cmd, pending)
	return nil
}
