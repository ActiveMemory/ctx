//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
	"github.com/ActiveMemory/ctx/internal/config/file"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/config"
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
	if target == file.ProfileProd {
		target = file.ProfileBase
	}

	var profile string
	switch target {
	case file.ProfileDev:
		profile = file.ProfileDev
	case file.ProfileBase:
		profile = file.ProfileBase
	case "":
		// Toggle.
		current := core.DetectProfile()
		if current == file.ProfileDev {
			profile = file.ProfileBase
		} else {
			profile = file.ProfileDev
		}
	default:
		return ctxErr.UnknownProfile(target)
	}

	msg, switchErr := core.SwitchTo(root, profile)
	if switchErr != nil {
		return switchErr
	}
	cmd.Println(msg)
	return nil
}
