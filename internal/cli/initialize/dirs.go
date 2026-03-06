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

// projectDirs lists the project-root directories created by ctx init,
// each with an explanatory README.md.
var projectDirs = []string{
	config.DirSpecs,
	config.DirIdeas,
}

// createProjectDirs creates project-root directories (specs/, ideas/) with
// README.md files. Skips directories that already exist. Creates the README
// inside new directories only.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func createProjectDirs(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()

	for _, dir := range projectDirs {
		if _, statErr := os.Stat(dir); statErr == nil {
			cmd.Println(fmt.Sprintf("  %s %s/ (exists, skipped)",
				color.YellowString("○"), dir))
			continue
		}

		if mkdirErr := os.MkdirAll(dir, config.PermExec); mkdirErr != nil {
			return fmt.Errorf("failed to create %s/: %w", dir, mkdirErr)
		}

		readme, readErr := assets.ProjectReadme(dir)
		if readErr != nil {
			return fmt.Errorf("failed to read %s README template: %w",
				dir, readErr)
		}

		readmePath := filepath.Join(dir, config.FilenameReadme)
		if writeErr := os.WriteFile(readmePath, readme, config.PermFile); writeErr != nil {
			return fmt.Errorf("failed to write %s: %w", readmePath, writeErr)
		}

		cmd.Println(fmt.Sprintf("  %s %s/", green("✓"), dir))
	}

	return nil
}
