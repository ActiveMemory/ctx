//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/config/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
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
	case file.ProfileDev:
		cmd.Println(desc.TextDesc(text.DescKeyWriteConfigProfileDev))
	case file.ProfileBase:
		cmd.Println(desc.TextDesc(text.DescKeyWriteConfigProfileBase))
	default:
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.DescKeyWriteConfigProfileNone),
			file.CtxRC,
		))
	}
	return nil
}
