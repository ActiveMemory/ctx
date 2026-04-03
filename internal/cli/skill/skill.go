//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/skill/cmd/install"
	"github.com/ActiveMemory/ctx/internal/cli/skill/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/skill/cmd/remove"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx skill" parent command.
//
// Returns:
//   - *cobra.Command: Configured skill command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeySkill, cmd.UseSkill,
		install.Cmd(),
		list.Cmd(),
		remove.Cmd(),
	)
}
