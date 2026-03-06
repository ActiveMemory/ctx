//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
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
			return runList(cmd)
		},
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(showCmd())
	cmd.AddCommand(addCmd())
	cmd.AddCommand(rmCmd())

	return cmd
}

// promptsDir returns the path to the prompts directory.
func promptsDir() string {
	return filepath.Join(rc.ContextDir(), config.DirPrompts)
}

// listCmd returns the prompt list subcommand.
func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available prompt templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd)
		},
	}
}

// runList prints all available prompt template names.
func runList(cmd *cobra.Command) error {
	dir := promptsDir()

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			cmd.Println("No prompts found. Run 'ctx init' or 'ctx prompt add' to create prompts.")
			return nil
		}
		return fmt.Errorf("read prompts directory: %w", err)
	}

	var found bool
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, config.ExtMarkdown) {
			continue
		}
		cmd.Println(fmt.Sprintf("  %s", strings.TrimSuffix(name, config.ExtMarkdown)))
		found = true
	}

	if !found {
		cmd.Println("No prompts found. Run 'ctx init' or 'ctx prompt add' to create prompts.")
	}

	return nil
}
