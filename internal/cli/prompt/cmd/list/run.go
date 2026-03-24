//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	"github.com/ActiveMemory/ctx/internal/config/file"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/prompt"
)

// Run prints all available prompt template names.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read failure
func Run(cmd *cobra.Command) error {
	dir := core.PromptsDir()

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			prompt.PromptNone(cmd)
			return nil
		}
		return errFs.ReadDirectory(dir, readErr)
	}

	var found bool
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, file.ExtMarkdown) {
			continue
		}
		prompt.PromptItem(cmd, strings.TrimSuffix(name, file.ExtMarkdown))
		found = true
	}

	if !found {
		prompt.PromptNone(cmd)
	}

	return nil
}
