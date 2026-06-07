//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reject

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/dream/core/dispose"
)

// Run delegates to the dispose core Reject logic.
//
// Parameters:
//   - cmd: cobra command for output
//   - id: the proposal ID to reject
//   - note: optional human note
//
// Returns:
//   - error: non-nil on a resolution, not-found, or ledger failure
func Run(cmd *cobra.Command, id, note string) error {
	return dispose.Reject(cmd, id, note)
}
