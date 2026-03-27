//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lock

import (
	coreLock "github.com/ActiveMemory/ctx/internal/cli/journal/core/lock"
	"github.com/spf13/cobra"
)

// runLockUnlock delegates to core.RunLockUnlock with lock=true.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Patterns to match against journal filenames
//   - all: If true, apply to all journal entries
//   - lock: True for lock, false for unlock
//
// Returns:
//   - error: Non-nil on validation or I/O failure
func runLockUnlock(
	cmd *cobra.Command,
	args []string,
	all, lock bool,
) error {
	return coreLock.RunLockUnlock(cmd, args, all, lock)
}
