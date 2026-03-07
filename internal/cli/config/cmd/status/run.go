//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run prints the active .ctxrc profile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - root: Git repository root directory
//
// Returns:
//   - error: Always nil (included for RunE compatibility)
func Run(cmd *cobra.Command, root string) error {
	profile := core.DetectProfile(root)
	switch profile {
	case core.ProfileDev:
		write.InfoConfigProfileDev(cmd)
	case core.ProfileBase:
		write.InfoConfigProfileBase(cmd)
	default:
		write.InfoConfigProfileNone(cmd, core.FileCtxRC)
	}
	return nil
}
