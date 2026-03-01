//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx guide" cobra command.
func Cmd() *cobra.Command {
	var (
		showSkills   bool
		showCommands bool
	)

	cmd := &cobra.Command{
		Use:   "guide",
		Short: "Quick-reference cheat sheet for ctx",
		Long: `Use-case-oriented cheat sheet for ctx.

Shows core commands grouped by workflow, key skills, and common recipes.
Default output fits one screen.

Use --skills to list all available slash-command skills.
Use --commands to list all CLI commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGuide(cmd, showSkills, showCommands)
		},
	}

	cmd.Flags().BoolVar(&showSkills, "skills", false, "List all available skills")
	cmd.Flags().BoolVar(&showCommands, "commands", false, "List all CLI commands")

	return cmd
}
