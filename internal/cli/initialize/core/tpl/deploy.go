//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// DeployTemplates creates a subdirectory under contextDir and writes embedded
// templates into it, skipping files that already exist unless force is true.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//   - force: If true, overwrite existing files
//   - p: Deploy parameters (subdirectory and error text keys)
//   - list: Returns the names of embedded templates
//   - read: Returns the content of a named template
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func DeployTemplates(
	cmd *cobra.Command, contextDir string, force bool,
	p entity.DeployParams,
	list func() ([]string, error),
	read func(string) ([]byte, error),
) error {
	targetDir := filepath.Join(contextDir, p.SubDir)
	if mkdirErr := os.MkdirAll(targetDir, fs.PermExec); mkdirErr != nil {
		return errFs.Mkdir(targetDir, mkdirErr)
	}

	names, listErr := list()
	if listErr != nil {
		return fmt.Errorf(desc.Text(p.ListErrKey), listErr)
	}

	for _, name := range names {
		targetPath := filepath.Join(targetDir, name)
		if _, statErr := os.Stat(targetPath); statErr == nil && !force {
			initialize.Skipped(cmd, filepath.Join(p.SubDir, name))
			continue
		}

		content, readErr := read(name)
		if readErr != nil {
			return fmt.Errorf(desc.Text(p.ReadErrKey), name, readErr)
		}

		if writeErr := os.WriteFile(
			targetPath, content, fs.PermFile,
		); writeErr != nil {
			return errFs.FileWrite(targetPath, writeErr)
		}

		initialize.Created(cmd, filepath.Join(p.SubDir, name))
	}

	return nil
}
