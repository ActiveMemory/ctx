//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/tpl"
)

// Commands returns the list of embedded command file names.
//
// These are Claude Code slash command definitions (e.g., "ctx-status.md",
// "ctx-reflect.md") from internal/templates/claude/commands/. They can be
// installed to .claude/commands/ via "ctx init".
//
// Returns:
//   - []string: Filenames of available command definitions
//   - error: Non-nil if the commands directory cannot be read
func Commands() ([]string, error) {
	names, err := tpl.ListClaudeCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to list commands: %w", err)
	}
	return names, nil
}

// CommandByName returns the content of a command file by name.
//
// Parameters:
//   - name: Filename as returned by [Commands] (e.g., "ctx-status.md")
//
// Returns:
//   - []byte: Raw bytes of the command definition file
//   - error: Non-nil if the command file does not exist or cannot be read
func CommandByName(name string) ([]byte, error) {
	content, err := tpl.ClaudeCommandByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read command %s: %w", name, err)
	}
	return content, nil
}
