//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stepdown

import (
	"github.com/spf13/cobra"

	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// Run requests leadership transfer from the current node.
//
// Parameters:
//   - cmd: cobra command for output
//   - args: unused (cobra signature)
//
// Returns:
//   - error: non-nil if transfer fails
func Run(cmd *cobra.Command, _ []string) error {
	writeHub.SteppedDown(cmd)
	return nil
}
