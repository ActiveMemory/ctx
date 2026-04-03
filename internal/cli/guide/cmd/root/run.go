//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/guide/core/command"
	"github.com/ActiveMemory/ctx/internal/cli/guide/core/skill"
	"github.com/ActiveMemory/ctx/internal/write/guide"
)

// Run dispatches to the appropriate guide output based on
// flags.
//
// Parameters:
//   - cmd: Cobra command for output stream and root traversal
//   - showSkills: If true, list all available skills
//   - showCommands: If true, list all CLI commands
//
// Returns:
//   - error: Non-nil if skill listing fails
func Run(
	cmd *cobra.Command, showSkills, showCommands bool,
) error {
	switch {
	case showSkills:
		return skill.List(cmd)
	case showCommands:
		return command.List(cmd)
	default:
		guide.Default(cmd)
		return nil
	}
}
