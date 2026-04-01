//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	writeVscode "github.com/ActiveMemory/ctx/internal/write/vscode"
)

// CreateVSCodeArtifacts generates VS Code-native configuration files
// as the editor-specific counterpart to Claude Code's settings and hooks.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation fails (file-level errors
//     are non-fatal and reported inline)
func CreateVSCodeArtifacts(cmd *cobra.Command) error {
	if mkdirErr := os.MkdirAll(cfgVscode.Dir, fs.PermExec); mkdirErr != nil {
		return mkdirErr
	}

	if extErr := writeExtensionsJSON(cmd); extErr != nil {
		writeVscode.InfoWarnNonFatal(cmd, cfgVscode.FileExtensionsJSON, extErr)
	}

	if taskErr := writeTasksJSON(cmd); taskErr != nil {
		writeVscode.InfoWarnNonFatal(cmd, cfgVscode.FileTasksJSON, taskErr)
	}

	if mcpErr := writeMCPJSON(cmd); mcpErr != nil {
		writeVscode.InfoWarnNonFatal(cmd, cfgVscode.FileMCPJSON, mcpErr)
	}

	return nil
}
