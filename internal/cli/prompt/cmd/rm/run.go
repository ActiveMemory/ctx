//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/file"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/prompt"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
)

// Run deletes a prompt template by name.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Template name (without .md extension)
//
// Returns:
//   - error: Non-nil on missing template or remove failure
func Run(cmd *cobra.Command, name string) error {
	path := filepath.Join(core.PromptsDir(), name+file.ExtMarkdown)

	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return ctxerr.NotFound(name)
	}

	if removeErr := os.Remove(path); removeErr != nil {
		return ctxerr.Remove(removeErr)
	}

	prompt.PromptRemoved(cmd, name)
	return nil
}
