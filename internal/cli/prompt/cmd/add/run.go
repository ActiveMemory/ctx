//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/prompt/core"
	"github.com/ActiveMemory/ctx/internal/config"
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
	if mkdirErr := os.MkdirAll(dir, config.PermExec); mkdirErr != nil {
		return fmt.Errorf("create prompts directory: %w", mkdirErr)
	}

	path := filepath.Join(dir, name+config.ExtMarkdown)

	// Check if file already exists.
	if _, statErr := os.Stat(path); statErr == nil {
		return fmt.Errorf("prompt %q already exists", name)
	}

	var content []byte

	if fromStdin {
		var readErr error
		content, readErr = io.ReadAll(cmd.InOrStdin())
		if readErr != nil {
			return fmt.Errorf("read stdin: %w", readErr)
		}
	} else {
		// Try to load from embedded starter templates.
		var templateErr error
		content, templateErr = assets.PromptTemplate(name + config.ExtMarkdown)
		if templateErr != nil {
			return fmt.Errorf("no embedded template %q — use --stdin to provide content", name)
		}
	}

	if writeErr := os.WriteFile(path, content, config.PermFile); writeErr != nil {
		return fmt.Errorf("write prompt: %w", writeErr)
	}

	cmd.Println(fmt.Sprintf("Created prompt %q.", name))
	return nil
}
