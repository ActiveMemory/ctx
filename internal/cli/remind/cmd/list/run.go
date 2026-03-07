//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind/core"
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
		cmd.Println("No reminders.")
		return nil
	}

	today := time.Now().Format("2006-01-02")
	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf("  (after %s, not yet due)", *r.After)
			}
		}
		cmd.Println(fmt.Sprintf("  [%d] %s%s", r.ID, r.Message, annotation))
	}

	return nil
}
