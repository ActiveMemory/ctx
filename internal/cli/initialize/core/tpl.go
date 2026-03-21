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

	"github.com/ActiveMemory/ctx/internal/assets/read/entry"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
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
func CreateEntryTemplates(
	cmd *cobra.Command, contextDir string, force bool,
) error {
	templatesDir := filepath.Join(contextDir, dir.Templates)
	if err := os.MkdirAll(templatesDir, fs.PermExec); err != nil {
		return errFs.Mkdir(templatesDir, err)
	}
	entryTemplates, err := entry.List()
	if err != nil {
		return ctxerr.ListEntryTemplates(err)
	}
	for _, name := range entryTemplates {
		targetPath := filepath.Join(templatesDir, name)
		if _, err := os.Stat(targetPath); err == nil && !force {
			initialize.Skipped(cmd, filepath.Join(dir.Templates, name))
			continue
		}
		content, err := entry.ForName(name)
		if err != nil {
			return ctxerr.ReadEntryTemplate(name, err)
		}
		if err := os.WriteFile(targetPath, content, fs.PermFile); err != nil {
			return errFs.FileWrite(targetPath, err)
		}
		initialize.Created(cmd, filepath.Join(dir.Templates, name))
	}
	return nil
}
