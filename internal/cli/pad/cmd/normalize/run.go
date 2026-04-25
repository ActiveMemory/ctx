//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/rc"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run reassigns entry IDs as 1..N in current file order.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write failure
func Run(cmd *cobra.Command) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	entries, readErr := store.ReadEntriesWithIDs()
	if readErr != nil {
		return readErr
	}

	if len(entries) == 0 {
		writePad.Empty(cmd)
		return nil
	}

	normalized := parse.Normalize(entries)

	writeErr := store.WriteEntriesWithIDs(cmd, normalized)
	if writeErr != nil {
		return writeErr
	}

	writePad.Normalized(cmd, len(normalized))
	return nil
}
