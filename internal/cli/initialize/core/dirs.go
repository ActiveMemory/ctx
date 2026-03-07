//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// ProjectDirs lists the project-root directories created by ctx init,
// each with an explanatory README.md.
var ProjectDirs = []string{
	config.DirSpecs,
	config.DirIdeas,
}

// CreateProjectDirs creates project-root directories (specs/, ideas/) with
// README.md files. Skips directories that already exist. Creates the README
// inside new directories only.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func CreateProjectDirs(cmd *cobra.Command) error {
	for _, dir := range ProjectDirs {
		if _, statErr := os.Stat(dir); statErr == nil {
			write.InitSkippedDir(cmd, dir)
			continue
		}

		if mkdirErr := os.MkdirAll(dir, config.PermExec); mkdirErr != nil {
			return ctxerr.Mkdir(dir+"/", mkdirErr)
		}

		readme, readErr := assets.ProjectReadme(dir)
		if readErr != nil {
			return ctxerr.ReadProjectReadme(dir, readErr)
		}

		readmePath := filepath.Join(dir, config.FilenameReadme)
		if writeErr := os.WriteFile(readmePath, readme, config.PermFile); writeErr != nil {
			return ctxerr.FileWrite(readmePath, writeErr)
		}

		write.InitCreatedDir(cmd, dir)
	}

	return nil
}
