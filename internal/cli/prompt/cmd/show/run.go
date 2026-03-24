//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/err/fs"
	errPrompt "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/io"
	writePrompt "github.com/ActiveMemory/ctx/internal/write/prompt"
)

// Run reads and prints a prompt template by name.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Template name (without .md extension)
//
// Returns:
//   - error: Non-nil on read failure or missing template
func Run(cmd *cobra.Command, name string) error {
	content, readErr := io.SafeReadFile(
		core.PromptsDir(), name+file.ExtMarkdown,
	)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return errPrompt.NotFound(name)
		}
		return fs.ReadFile(readErr)
	}

	writePrompt.ShowContent(cmd, content)
	return nil
}
