//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/date"
	"github.com/ActiveMemory/ctx/internal/write/remind"
)

// Run creates a new reminder and prints confirmation.
//
// Exported for reuse by the parent command's default action.
//
// Parameters:
//   - cmd: Cobra command for output
//   - message: Reminder text
//   - after: Optional date gate in YYYY-MM-DD format (empty string to skip)
//
// Returns:
//   - error: Non-nil on read/write failure or invalid date
func Run(cmd *cobra.Command, message, after string) error {
	reminders, readErr := core.ReadReminders()
	if readErr != nil {
		return readErr
	}

	r := core.Reminder{
		ID:      core.NextID(reminders),
		Message: message,
		Created: time.Now().UTC().Format(time.RFC3339),
	}
	if after != "" {
		if _, parseErr := time.Parse(cfgTime.DateFormat, after); parseErr != nil {
			return ctxerr.InvalidDateValue(after)
		}
		r.After = &after
	}

	reminders = append(reminders, r)
	if writeErr := core.WriteReminders(reminders); writeErr != nil {
		return writeErr
	}

	remind.ReminderAdded(cmd, r.ID, r.Message, r.After)
	return nil
}
