//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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

	short, long := desc.Command(cmd.DescKeyGuide)
	c := &cobra.Command{
		Use:         cmd.UseGuide,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: ""},
		Long:        long,
		Example:     desc.Example(cmd.DescKeyGuide),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, showSkills, showCommands)
		},
	}

	flagbind.BoolFlag(
		c, &showSkills,
		cFlag.Skills, flag.DescKeyGuideSkills,
	)
	flagbind.BoolFlag(
		c, &showCommands,
		cFlag.Commands, flag.DescKeyGuideCommands,
	)

	return c
}
