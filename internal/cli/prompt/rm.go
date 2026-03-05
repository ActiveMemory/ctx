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

// rmCmd returns the prompt rm subcommand.
//
// Returns:
//   - *cobra.Command: Configured rm subcommand
func rmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm NAME",
		Short: "Remove a prompt template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRm(cmd, args[0])
		},
	}
}

// runRm deletes a prompt template by name.
func runRm(cmd *cobra.Command, name string) error {
	path := filepath.Join(promptsDir(), name+config.ExtMarkdown)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("prompt %q not found", name)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove prompt: %w", err)
	}

	cmd.Println(fmt.Sprintf("Removed prompt %q.", name))
	return nil
}
