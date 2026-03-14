//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx guide" cobra command.
//
// Returns:
//   - *cobra.Command: Configured guide command with flags registered
func Cmd() *cobra.Command {
	var (
		showSkills   bool
		showCommands bool
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyGuide)
	cmd := &cobra.Command{
		Use:         "guide",
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: ""},
		Long:        long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, showSkills, showCommands)
		},
	}

	cmd.Flags().BoolVar(
		&showSkills,
		"skills",
		false,
		assets.FlagDesc(assets.FlagDescKeyGuideSkills),
	)
	cmd.Flags().BoolVar(
		&showCommands,
		"commands",
		false,
		assets.FlagDesc(assets.FlagDescKeyGuideCommands),
	)

	return cmd
}
