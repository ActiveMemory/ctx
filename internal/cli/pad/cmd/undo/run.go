//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package undo

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/rc"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run restores the scratchpad from the most recent on-disk
// snapshot. Snapshots are written transparently before every
// destructive pad mutation; see [store.SnapshotBefore].
//
// Empty history is not an error: prints a friendly message and
// exits 0 so cron'd or scripted invocations stay quiet.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: non-nil only on history-read, snapshot-take, or
//     restore-copy failures
func Run(cmd *cobra.Command) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}

	slot, restoreErr := store.Restore(cmd)
	if restoreErr != nil {
		return restoreErr
	}

	if slot == "" {
		writePad.NoHistory(cmd)
		return nil
	}

	writePad.Restored(cmd, slot)
	return nil
}
