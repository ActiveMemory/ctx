//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/validate"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run removes entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//
// Returns:
//   - error: Non-nil on invalid index or read/write failure
func Run(cmd *cobra.Command, n int) error {
	entries, err := store.ReadEntries()
	if err != nil {
		return err
	}

	if validErr := validate.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}

	entries = append(entries[:n-1], entries[n:]...)

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	pad.EntryRemoved(cmd, n)
	return nil
}
