//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unlock

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/lock"
)

// Run delegates to core.Run with lock=false.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Patterns to match against journal filenames
//   - all: If true, apply to all journal entries
//
// Returns:
//   - error: Non-nil on validation or I/O failure
func Run(
	cmd *cobra.Command,
	args []string,
	all bool,
) error {
	return lock.Run(cmd, args, all, false)
}
