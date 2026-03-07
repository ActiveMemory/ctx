//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core"
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
		if _, parseErr := time.Parse("2006-01-02", after); parseErr != nil {
			return fmt.Errorf("invalid date %q (expected YYYY-MM-DD)", after)
		}
		r.After = &after
	}

	reminders = append(reminders, r)
	if writeErr := core.WriteReminders(reminders); writeErr != nil {
		return writeErr
	}

	suffix := ""
	if r.After != nil {
		suffix = fmt.Sprintf("  (after %s)", *r.After)
	}
	cmd.Println(fmt.Sprintf("  + [%d] %s%s", r.ID, r.Message, suffix))
	return nil
}
