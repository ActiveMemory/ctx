//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package vscode

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	writeVscode "github.com/ActiveMemory/ctx/internal/write/vscode"
)

// CreateVSCodeArtifacts generates VS Code workspace configuration files
// (.vscode/) during ctx init.
//
// Creates extensions.json, tasks.json, and mcp.json as the
// editor-specific counterpart to Claude Code's settings and hooks.
// Individual file errors are non-fatal and reported inline.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation fails
func CreateVSCodeArtifacts(cmd *cobra.Command) error {
	if mkdirErr := os.MkdirAll(cfgVscode.Dir, fs.PermExec); mkdirErr != nil {
		return mkdirErr
	}

	if extErr := createExtensionsJSON(cmd); extErr != nil {
		writeVscode.InfoWarnNonFatal(cmd, cfgVscode.FileExtensionsJSON, extErr)
	}

	if taskErr := createTasksJSON(cmd); taskErr != nil {
		writeVscode.InfoWarnNonFatal(cmd, cfgVscode.FileTasksJSON, taskErr)
	}

	if mcpErr := createMCPJSON(cmd); mcpErr != nil {
		writeVscode.InfoWarnNonFatal(cmd, cfgVscode.FileMCPJSON, mcpErr)
	}

	return nil
}
