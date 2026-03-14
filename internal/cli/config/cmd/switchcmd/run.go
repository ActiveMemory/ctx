//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// Run executes the profile switch logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - root: Git repository root directory
//   - args: Optional profile name (dev, base, or prod alias)
//
// Returns:
//   - error: Non-nil on unknown profile or copy failure
func Run(cmd *cobra.Command, root string, args []string) error {
	var target string
	if len(args) > 0 {
		target = args[0]
	}

	// Normalize "prod" alias.
	if target == core.ProfileProd {
		target = core.ProfileBase
	}

	var profile string
	switch target {
	case core.ProfileDev:
		profile = core.ProfileDev
	case core.ProfileBase:
		profile = core.ProfileBase
	case "":
		// Toggle.
		current := core.DetectProfile()
		if current == core.ProfileDev {
			profile = core.ProfileBase
		} else {
			profile = core.ProfileDev
		}
	default:
		return ctxerr.UnknownProfile(target)
	}

	msg, switchErr := core.SwitchTo(root, profile)
	if switchErr != nil {
		return switchErr
	}
	cmd.Println(msg)
	return nil
}
