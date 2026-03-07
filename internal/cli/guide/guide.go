//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	guideroot "github.com/ActiveMemory/ctx/internal/cli/guide/cmd/root"
	"github.com/ActiveMemory/ctx/internal/config"
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

	short, long := assets.CommandDesc("guide")
	cmd := &cobra.Command{
		Use:         "guide",
		Short:       short,
		Annotations: map[string]string{config.AnnotationSkipInit: ""},
		Long:        long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return guideroot.Run(cmd, showSkills, showCommands)
		},
	}

	cmd.Flags().BoolVar(&showSkills, "skills", false, "List all available skills")
	cmd.Flags().BoolVar(&showCommands, "commands", false, "List all CLI commands")

	return cmd
}
