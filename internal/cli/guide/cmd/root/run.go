//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Run dispatches to the appropriate guide output based on flags.
//
// Parameters:
//   - cmd: Cobra command for output stream and root traversal
//   - showSkills: If true, list all available skills
//   - showCommands: If true, list all CLI commands
//
// Returns:
//   - error: Non-nil if skill listing fails
func Run(cmd *cobra.Command, showSkills, showCommands bool) error {
	switch {
	case showSkills:
		return listSkills(cmd)
	case showCommands:
		return listCommands(cmd)
	default:
		cmd.Print(desc.Text(text.DescKeyGuideDefault))
		return nil
	}
}
