//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"fmt"
	"github.com/ActiveMemory/ctx/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// showCmd returns the prompt show subcommand.
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show NAME",
		Short: "Print a prompt template to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(cmd, args[0])
		},
	}
}

// runShow reads and prints a prompt template by name.
func runShow(cmd *cobra.Command, name string) error {
	path := filepath.Join(promptsDir(), name+config.ExtMarkdown)

	content, err := os.ReadFile(path) //nolint:gosec // user-provided name is intentional
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("prompt %q not found", name)
		}
		return fmt.Errorf("read prompt: %w", err)
	}

	cmd.Print(string(content))
	return nil
}
