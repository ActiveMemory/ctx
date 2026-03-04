//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// addCmd returns the prompt add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func addCmd() *cobra.Command {
	var fromStdin bool

	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: "Create a new prompt from embedded template or stdin",
		Long: `Create a new prompt template in .context/prompts/.

By default, creates from an embedded starter template if one exists
with the given name. Use --stdin to read content from standard input.

Examples:
  ctx prompt add code-review
  echo "# My Prompt" | ctx prompt add my-prompt --stdin`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0], fromStdin)
		},
	}

	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "read prompt content from stdin")

	return cmd
}

// runAdd creates a new prompt template file.
func runAdd(cmd *cobra.Command, name string, fromStdin bool) error {
	dir := promptsDir()
	if err := os.MkdirAll(dir, config.PermExec); err != nil {
		return fmt.Errorf("create prompts directory: %w", err)
	}

	path := filepath.Join(dir, name+".md")

	// Check if file already exists.
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("prompt %q already exists", name)
	}

	var content []byte

	if fromStdin {
		var err error
		content, err = io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
	} else {
		// Try to load from embedded starter templates.
		var err error
		content, err = assets.PromptTemplate(name + ".md")
		if err != nil {
			return fmt.Errorf("no embedded template %q — use --stdin to provide content", name)
		}
	}

	if err := os.WriteFile(path, content, config.PermFile); err != nil {
		return fmt.Errorf("write prompt: %w", err)
	}

	cmd.Printf("Created prompt %q.\n", name)
	return nil
}
