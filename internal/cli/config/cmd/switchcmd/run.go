//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core/profile"
	"github.com/ActiveMemory/ctx/internal/config/file"
	errConfig "github.com/ActiveMemory/ctx/internal/err/config"
	writeConfig "github.com/ActiveMemory/ctx/internal/write/config"
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

	var p string
	switch target {
	case file.ProfileDev:
		p = file.ProfileDev
	case file.ProfileBase:
		p = file.ProfileBase
	case "":
		// Toggle.
		current := profile.Detect()
		if current == file.ProfileDev {
			p = file.ProfileBase
		} else {
			p = file.ProfileDev
		}
	default:
		return errConfig.UnknownProfile(target)
	}

	msg, switchErr := profile.SwitchTo(root, p)
	if switchErr != nil {
		return switchErr
	}
	writeConfig.SwitchConfirm(cmd, msg)
	return nil
}
