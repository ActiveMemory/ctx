//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/rm"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/cmd/show"
)

// Cmd returns the prompt command with subcommands.
//
// When invoked without a subcommand, it lists all prompt templates.
//
// Returns:
//   - *cobra.Command: Configured prompt command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Manage reusable prompt templates",
		Long: `Manage prompt templates stored in .context/prompts/.

Prompt templates are plain markdown files — no frontmatter, no build step.
Use them as lightweight, reusable instructions for common tasks like
code reviews, refactoring, or explaining code.

When invoked without a subcommand, lists all available prompts.

Subcommands:
  list     List available prompt templates
  show     Print a prompt template to stdout
  add      Create a new prompt from embedded template or stdin
  rm       Remove a prompt template`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return list.RunList(cmd)
		},
	}

	cmd.AddCommand(list.Cmd())
	cmd.AddCommand(show.Cmd())
	cmd.AddCommand(add.Cmd())
	cmd.AddCommand(rm.Cmd())

	return cmd
}
