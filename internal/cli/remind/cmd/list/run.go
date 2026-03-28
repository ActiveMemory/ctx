//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/write/remind"
)

// Run prints all pending reminders with date annotations.
//
// Exported for reuse by the parent command's default action.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read failure
func Run(cmd *cobra.Command) error {
	reminders, readErr := core.ReadReminders()
	if readErr != nil {
		return readErr
	}

	if len(reminders) == 0 {
		remind.None(cmd)
		return nil
	}

	today := time.Now().Format(cfgTime.DateFormat)
	for _, r := range reminders {
		remind.Item(cmd, r.ID, r.Message, r.After, today)
	}

	return nil
}
