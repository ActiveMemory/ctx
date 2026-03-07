//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// CreatePromptTemplates creates prompt template files in .context/prompts/.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//   - force: If true, overwrite existing files
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func CreatePromptTemplates(cmd *cobra.Command, contextDir string, force bool) error {
	promptDir := filepath.Join(contextDir, config.DirPrompts)
	if err := os.MkdirAll(promptDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", promptDir, err)
	}
	promptTemplates, err := assets.ListPromptTemplates()
	if err != nil {
		return fmt.Errorf("failed to list prompt templates: %w", err)
	}
	for _, name := range promptTemplates {
		targetPath := filepath.Join(promptDir, name)
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Println(fmt.Sprintf("  ○ prompts/%s (exists, skipped)", name))
			continue
		}
		content, err := assets.PromptTemplate(name)
		if err != nil {
			return fmt.Errorf("failed to read prompt template %s: %w", name, err)
		}
		if err := os.WriteFile(targetPath, content, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}
		cmd.Println(fmt.Sprintf("  ✓ prompts/%s", name))
	}
	return nil
}
