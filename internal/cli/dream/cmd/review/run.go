//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package review

import (
	"github.com/spf13/cobra"

	coreReview "github.com/ActiveMemory/ctx/internal/cli/dream/core/review"
)

// Run delegates to the review core logic.
//
// Parameters:
//   - cmd: cobra command for output
//
// Returns:
//   - error: non-nil on a resolution or read failure
func Run(cmd *cobra.Command) error {
	return coreReview.Run(cmd)
}
