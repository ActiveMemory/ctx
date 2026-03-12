//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
func CreatePromptTemplates(cmd *cobra.Command, contextDir string, force bool) error {
	promptDir := filepath.Join(contextDir, dir.Prompts)
	if err := os.MkdirAll(promptDir, fs.PermExec); err != nil {
		return ctxerr.Mkdir(promptDir, err)
	}
	promptTemplates, err := assets.ListPromptTemplates()
	if err != nil {
		return ctxerr.ListPromptTemplates(err)
	}
	for _, name := range promptTemplates {
		targetPath := filepath.Join(promptDir, name)
		if _, err := os.Stat(targetPath); err == nil && !force {
			write.InitSkipped(cmd, "prompts/"+name)
			continue
		}
		content, err := assets.PromptTemplate(name)
		if err != nil {
			return ctxerr.ReadPromptTemplate(name, err)
		}
		if err := os.WriteFile(targetPath, content, fs.PermFile); err != nil {
			return ctxerr.FileWrite(targetPath, err)
		}
		write.InitCreated(cmd, "prompts/"+name)
	}
	return nil
}
