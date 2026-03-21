//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/project"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"
)

// ProjectDirs lists the project-root directories created by ctx init,
// each with an explanatory README.md.
var ProjectDirs = []string{
	dir.Specs,
	dir.Ideas,
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
			initialize.SkippedDir(cmd, dir)
			continue
		}

		if mkdirErr := os.MkdirAll(dir, fs.PermExec); mkdirErr != nil {
			return fs2.Mkdir(dir+"/", mkdirErr)
		}

		readme, readErr := project.Readme(dir)
		if readErr != nil {
			return ctxerr.ReadProjectReadme(dir, readErr)
		}

		readmePath := filepath.Join(dir, file.Readme)
		if writeErr := os.WriteFile(readmePath, readme, fs.PermFile); writeErr != nil {
			return fs2.FileWrite(readmePath, writeErr)
		}

		initialize.CreatedDir(cmd, dir)
	}

	return nil
}
