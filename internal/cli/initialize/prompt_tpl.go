//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// createPromptTemplates creates .context/prompts/ with starter prompt
// templates for common tasks like code review and refactoring.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - contextDir: Path to the .context/ directory
//   - force: If true, overwrite existing templates
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func createPromptTemplates(
	cmd *cobra.Command, contextDir string, force bool,
) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	promptDir := filepath.Join(contextDir, config.DirPrompts)
	if err := os.MkdirAll(promptDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", promptDir, err)
	}

	// Get list of prompt templates
	promptTemplates, err := assets.ListPromptTemplates()
	if err != nil {
		return fmt.Errorf("failed to list prompt templates: %w", err)
	}

	for _, name := range promptTemplates {
		targetPath := filepath.Join(promptDir, name)

		// Check if the file exists and --force not set
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Printf("  %s prompts/%s (exists, skipped)\n", yellow("○"), name)
			continue
		}

		content, err := assets.PromptTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to read prompt template %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		cmd.Printf("  %s prompts/%s\n", green("✓"), name)
	}

	return nil
}
