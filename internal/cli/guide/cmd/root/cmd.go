//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
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

	short, long := desc.CommandDesc(cmd.DescKeyGuide)
	cmd := &cobra.Command{
		Use:         cmd.DescKeyGuide,
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
		desc.FlagDesc(flag.DescKeyGuideSkills),
	)
	cmd.Flags().BoolVar(
		&showCommands,
		"commands",
		false,
		desc.FlagDesc(flag.DescKeyGuideCommands),
	)

	return cmd
}
