//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
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
	profile := core.DetectProfile()
	switch profile {
	case core.ProfileDev:
		cmd.Println(assets.TextDesc(assets.TextDescKeyWriteConfigProfileDev))
	case core.ProfileBase:
		cmd.Println(assets.TextDesc(assets.TextDescKeyWriteConfigProfileBase))
	default:
		cmd.Println(fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyWriteConfigProfileNone),
			core.FileCtxRC,
		))
	}
	return nil
}
