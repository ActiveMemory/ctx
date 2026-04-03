//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	coreDismiss "github.com/ActiveMemory/ctx/internal/cli/remind/core/dismiss"
)

// Run dismisses one or all reminders based on the
// all flag.
//
// When all is true, removes every reminder. Otherwise
// removes the single reminder identified by idStr.
//
// Parameters:
//   - cmd: Cobra command for output
//   - idStr: String reminder ID (ignored when all
//     is true)
//   - all: When true, dismiss all reminders
//
// Returns:
//   - error: Non-nil on invalid ID, missing reminder,
//     or write failure
func Run(
	cmd *cobra.Command, idStr string, all bool,
) error {
	if all {
		return coreDismiss.All(cmd)
	}
	return coreDismiss.One(cmd, idStr)
}
