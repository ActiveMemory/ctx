//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package amend

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/dream/core/dispose"
)

// Run delegates to the dispose core Amend logic.
//
// Parameters:
//   - cmd: cobra command for output
//   - id: the proposal ID to amend
//   - action: the action to apply instead of the recommendation
//   - note: optional human note
//
// Returns:
//   - error: non-nil on a resolution, not-found, unknown-action,
//     guard, mutation, or ledger failure
func Run(cmd *cobra.Command, id, action, note string) error {
	return dispose.Amend(cmd, id, action, note)
}
