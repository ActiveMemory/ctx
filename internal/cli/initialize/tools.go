//   /    Context:                     https://ctx.ist
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

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/tpl"
)

// createTools creates .context/tools/ with embedded tool scripts.
//
// Tool scripts are deployed as executable files that users can run
// directly from their project's .context/tools/ directory.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - contextDir: Path to the .context/ directory
//   - force: If true, overwrite existing tools
//
// Returns:
//   - error: Non-nil if directory creation or file operations fail
func createTools(cmd *cobra.Command, contextDir string, force bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	toolsDir := filepath.Join(contextDir, config.DirTools)
	if err := os.MkdirAll(toolsDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", toolsDir, err)
	}

	tools, err := tpl.ListTools()
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	for _, name := range tools {
		targetPath := filepath.Join(toolsDir, name)

		if _, err := os.Stat(targetPath); err == nil && !force {
			cmd.Printf("  %s tools/%s (exists, skipped)\n", yellow("○"), name)
			continue
		}

		content, err := tpl.Tool(name)
		if err != nil {
			return fmt.Errorf("failed to read tool %s: %w", name, err)
		}

		if err := os.WriteFile(targetPath, content, config.PermExec); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		cmd.Printf("  %s tools/%s\n", green("✓"), name)
	}

	return nil
}
