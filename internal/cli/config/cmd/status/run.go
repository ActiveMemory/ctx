//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
)

// RunStatus prints the active .ctxrc profile.
//
// Parameters:
//   - cmd: Cobra command for output
//   - root: Git repository root directory
//
// Returns:
//   - error: Always nil (included for RunE compatibility)
func RunStatus(cmd *cobra.Command, root string) error {
	profile := core.DetectProfile(root)
	switch profile {
	case core.ProfileDev:
		cmd.Println("active: dev (verbose logging enabled)")
	case core.ProfileBase:
		cmd.Println("active: base (defaults)")
	default:
		cmd.Println(fmt.Sprintf("active: none (%s does not exist)", core.FileCtxRC))
	}
	return nil
}
