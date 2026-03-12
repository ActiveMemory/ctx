//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
			write.PromptNone(cmd)
			return nil
		}
		return ctxerr.ReadDirectory(dir, readErr)
	}

	var found bool
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, file.ExtMarkdown) {
			continue
		}
		write.PromptItem(cmd, strings.TrimSuffix(name, file.ExtMarkdown))
		found = true
	}

	if !found {
		write.PromptNone(cmd)
	}

	return nil
}
