//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core/store"
	"github.com/ActiveMemory/ctx/internal/write/remind"
)

// Run reassigns reminder IDs as 1..N in current order.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write failure
func Run(cmd *cobra.Command) error {
	reminders, readErr := store.Read()
	if readErr != nil {
		return readErr
	}

	if len(reminders) == 0 {
		remind.None(cmd)
		return nil
	}

	for i := range reminders {
		reminders[i].ID = i + 1
	}

	if writeErr := store.Write(reminders); writeErr != nil {
		return writeErr
	}

	remind.Normalized(cmd, len(reminders))
	return nil
}
