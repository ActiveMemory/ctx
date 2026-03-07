//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	"github.com/ActiveMemory/ctx/internal/config"
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
			cmd.Println("No prompts found. Run 'ctx init' or 'ctx prompt add' to create prompts.")
			return nil
		}
		return fmt.Errorf("read prompts directory: %w", readErr)
	}

	var found bool
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, config.ExtMarkdown) {
			continue
		}
		cmd.Println(fmt.Sprintf("  %s", strings.TrimSuffix(name, config.ExtMarkdown)))
		found = true
	}

	if !found {
		cmd.Println("No prompts found. Run 'ctx init' or 'ctx prompt add' to create prompts.")
	}

	return nil
}
