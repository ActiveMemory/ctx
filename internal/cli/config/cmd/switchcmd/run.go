//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
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
	if target == "prod" {
		target = core.ProfileBase
	}

	switch target {
	case core.ProfileDev:
		return switchTo(cmd, root, core.ProfileDev)
	case core.ProfileBase:
		return switchTo(cmd, root, core.ProfileBase)
	case "":
		// Toggle.
		current := core.DetectProfile(root)
		if current == core.ProfileDev {
			return switchTo(cmd, root, core.ProfileBase)
		}
		return switchTo(cmd, root, core.ProfileDev)
	default:
		return fmt.Errorf(
			"unknown profile %q: must be dev, base, or prod", target)
	}
}

// switchTo copies the requested profile and prints a status message.
func switchTo(cmd *cobra.Command, root, profile string) error {
	current := core.DetectProfile(root)
	if current == profile {
		cmd.Println(fmt.Sprintf("already on %s profile", profile))
		return nil
	}

	var srcFile string
	if profile == core.ProfileDev {
		srcFile = core.FileCtxRCDev
	} else {
		srcFile = core.FileCtxRCBase
	}

	if copyErr := core.CopyProfile(root, srcFile); copyErr != nil {
		return copyErr
	}

	if current == "" {
		cmd.Println(fmt.Sprintf("created %s from %s profile", core.FileCtxRC, profile))
	} else {
		cmd.Println(fmt.Sprintf("switched to %s profile", profile))
	}
	return nil
}
