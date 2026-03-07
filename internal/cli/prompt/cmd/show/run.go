//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	"github.com/ActiveMemory/ctx/internal/config"
)

// runShow reads and prints a prompt template by name.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Template name (without .md extension)
//
// Returns:
//   - error: Non-nil on read failure or missing template
func runShow(cmd *cobra.Command, name string) error {
	path := filepath.Join(core.PromptsDir(), name+config.ExtMarkdown)

	content, readErr := os.ReadFile(path) //nolint:gosec // user-provided name is intentional
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return fmt.Errorf("prompt %q not found", name)
		}
		return fmt.Errorf("read prompt: %w", readErr)
	}

	cmd.Print(string(content))
	return nil
}
