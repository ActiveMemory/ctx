//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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
	if err := os.MkdirAll(templatesDir, fs.PermExec); err != nil {
		return fs2.Mkdir(templatesDir, err)
	}
	entryTemplates, err := assets.ListEntry()
	if err != nil {
		return ctxerr.ListEntryTemplates(err)
	}
	for _, name := range entryTemplates {
		targetPath := filepath.Join(templatesDir, name)
		if _, err := os.Stat(targetPath); err == nil && !force {
			initialize.Skipped(cmd, "templates/"+name)
			continue
		}
		content, err := assets.Entry(name)
		if err != nil {
			return ctxerr.ReadEntryTemplate(name, err)
		}
		if err := os.WriteFile(targetPath, content, fs.PermFile); err != nil {
			return fs2.FileWrite(targetPath, err)
		}
		initialize.Created(cmd, "templates/"+name)
	}
	return nil
}
