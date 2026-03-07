//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	"github.com/ActiveMemory/ctx/internal/config"
)

// runRm deletes a prompt template by name.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Template name (without .md extension)
//
// Returns:
//   - error: Non-nil on missing template or remove failure
func runRm(cmd *cobra.Command, name string) error {
	path := filepath.Join(core.PromptsDir(), name+config.ExtMarkdown)

	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return fmt.Errorf("prompt %q not found", name)
	}

	if removeErr := os.Remove(path); removeErr != nil {
		return fmt.Errorf("remove prompt: %w", removeErr)
	}

	cmd.Println(fmt.Sprintf("Removed prompt %q.", name))
	return nil
}
