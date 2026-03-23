//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/project"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errInitialize "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

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
	for _, d := range dirs {
		if _, statErr := os.Stat(d); statErr == nil {
			initialize.SkippedDir(cmd, d)
			continue
		}

		if mkdirErr := os.MkdirAll(d, fs.PermExec); mkdirErr != nil {
			return errFs.Mkdir(d, mkdirErr)
		}

		readme, readErr := project.Readme(d)
		if readErr != nil {
			return errInitialize.ReadProjectReadme(d, readErr)
		}

		readmePath := filepath.Join(d, file.Readme)
		if writeErr := os.WriteFile(readmePath, readme, fs.PermFile); writeErr != nil {
			return errFs.FileWrite(readmePath, writeErr)
		}

		initialize.CreatedDir(cmd, d)
	}

	return nil
}
