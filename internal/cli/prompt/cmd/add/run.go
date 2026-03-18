//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"io"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/prompt"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
)

// runAdd creates a new prompt template file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Template name (without .md extension)
//   - fromStdin: When true, read content from stdin instead of embedded templates
//
// Returns:
//   - error: Non-nil on write failure, duplicate name, or missing template
func runAdd(cmd *cobra.Command, name string, fromStdin bool) error {
	dir := core.PromptsDir()
	if mkdirErr := os.MkdirAll(dir, fs.PermExec); mkdirErr != nil {
		return fs2.Mkdir("prompts directory", mkdirErr)
	}

	path := filepath.Join(dir, name+file.ExtMarkdown)

	// Check if file already exists.
	if _, statErr := os.Stat(path); statErr == nil {
		return ctxerr.Exists(name)
	}

	var content []byte

	if fromStdin {
		var readErr error
		content, readErr = io.ReadAll(cmd.InOrStdin())
		if readErr != nil {
			return fs2.ReadInput(readErr)
		}
	} else {
		// Try to load from embedded starter templates.
		var templateErr error
		content, templateErr = assets.PromptTemplate(name + file.ExtMarkdown)
		if templateErr != nil {
			return ctxerr.NoPromptTemplate(name)
		}
	}

	if writeErr := os.WriteFile(path, content, fs.PermFile); writeErr != nil {
		return fs2.WriteFileFailed(writeErr)
	}

	prompt.PromptCreated(cmd, name)
	return nil
}
