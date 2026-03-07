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

// CreateEntryTemplates creates entry template files in .context/templates/.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//   - force: If true, overwrite existing files
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func CreateEntryTemplates(cmd *cobra.Command, contextDir string, force bool) error {
	templatesDir := filepath.Join(contextDir, "templates")
	if err := os.MkdirAll(templatesDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", templatesDir, err)
	}
	entryTemplates, err := assets.ListEntry()
	if err != nil {
		return fmt.Errorf("failed to list entry templates: %w", err)
	}
	for _, name := range entryTemplates {
		targetPath := filepath.Join(templatesDir, name)
		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Println(fmt.Sprintf("  ○ templates/%s (exists, skipped)", name))
			continue
		}
		content, err := assets.Entry(name)
		if err != nil {
			return fmt.Errorf("failed to read entry template %s: %w", name, err)
		}
		if err := os.WriteFile(targetPath, content, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}
		cmd.Println(fmt.Sprintf("  ✓ templates/%s", name))
	}
	return nil
}
