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

	"github.com/ActiveMemory/ctx/internal/assets/read/prompt"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
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
func CreatePromptTemplates(
	cmd *cobra.Command, contextDir string, force bool,
) error {
	promptDir := filepath.Join(contextDir, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		return errFs.Mkdir(promptDir, err)
	}
	promptTemplates, err := prompt.TemplateList()
	if err != nil {
		return ctxErr.ListPromptTemplates(err)
	}
	for _, name := range promptTemplates {
		targetPath := filepath.Join(promptDir, name)
		if _, err := os.Stat(targetPath); err == nil && !force {
			initialize.Skipped(cmd, filepath.Join(dir.Prompts, name))
			continue
		}
		content, err := prompt.Template(name)
		if err != nil {
			return ctxErr.ReadPromptTemplate(name, err)
		}
		if err := os.WriteFile(targetPath, content, fs.PermFile); err != nil {
			return errFs.FileWrite(targetPath, err)
		}
		initialize.Created(cmd, filepath.Join(dir.Prompts, name))
	}
	return nil
}
