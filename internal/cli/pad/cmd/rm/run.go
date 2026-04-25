//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run removes entries by stable ID. Supports multiple IDs.
// All IDs are resolved before any deletion to avoid
// shift-induced mismatches.
//
// Parameters:
//   - cmd: Cobra command for output
//   - ids: Stable entry IDs to remove
//
// Returns:
//   - error: Non-nil on invalid ID or read/write failure
func Run(cmd *cobra.Command, ids []int) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	entries, readErr := store.ReadEntriesWithIDs()
	if readErr != nil {
		return readErr
	}

	// Resolve all IDs before deleting any.
	removeSet := make(map[int]bool, len(ids))
	for _, id := range ids {
		idx := parse.FindByID(entries, id)
		if idx < 0 {
			return errPad.EntryNotFound(id)
		}
		removeSet[id] = true
	}

	// Filter out removed entries.
	var remaining []parse.Entry
	for _, e := range entries {
		if !removeSet[e.ID] {
			remaining = append(remaining, e)
		}
	}

	writeErr := store.WriteEntriesWithIDs(cmd, remaining)
	if writeErr != nil {
		return writeErr
	}

	for _, id := range ids {
		pad.EntryRemoved(cmd, id)
	}
	return nil
}
