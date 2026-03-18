//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mv

import (
	"github.com/ActiveMemory/ctx/internal/write/pad"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
)

// Run moves entry from 1-based position n to 1-based position m.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: Source position (1-based)
//   - m: Destination position (1-based)
//
// Returns:
//   - error: Non-nil on invalid index or read/write failure
func Run(cmd *cobra.Command, n, m int) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if validErr := core.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}
	if validErr := core.ValidateIndex(m, entries); validErr != nil {
		return validErr
	}

	// Extract the entry at position n
	entry := entries[n-1]
	// Remove it
	entries = append(entries[:n-1], entries[n:]...)
	// Insert at position m (adjust for 0-based)
	idx := m - 1
	entries = append(entries[:idx], append([]string{entry}, entries[idx:]...)...)

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	pad.PadEntryMoved(cmd, n, m)
	return nil
}
