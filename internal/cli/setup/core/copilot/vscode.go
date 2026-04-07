//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package copilot

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureVSCodeMCP creates .vscode/mcp.json to register the ctx MCP
// server for VS Code Copilot.
//
// Skips if the file already exists to preserve user customizations.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func ensureVSCodeMCP(cmd *cobra.Command) error {
	target := filepath.Join(cfgVscode.Dir, cfgVscode.FileMCPJSON)

	if _, statErr := os.Stat(target); statErr == nil {
		writeSetup.InfoCopilotCLISkipped(cmd, target)
		return nil
	}

	mkdirErr := ctxIo.SafeMkdirAll(
		cfgVscode.Dir, fs.PermExec,
	)
	if mkdirErr != nil {
		return mkdirErr
	}

	mcpCfg := map[string]interface{}{
		cfgVscode.KeyServers: map[string]interface{}{
			mcpServer.Name: map[string]interface{}{
				cfgVscode.KeyCommand: mcpServer.Command,
				cfgVscode.KeyArgs:    mcpServer.Args(),
			},
		},
	}
	data, _ := json.MarshalIndent(mcpCfg, "", token.Indent2)
	data = append(data, token.NewlineLF...)

	writeFileErr := ctxIo.SafeWriteFile(
		target, data, fs.PermFile,
	)
	if writeFileErr != nil {
		return writeFileErr
	}
	writeSetup.InfoCopilotCLICreated(cmd, target)
	return nil
}
