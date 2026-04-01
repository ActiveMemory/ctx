//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package copilotcli

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureMCPConfig registers the ctx MCP server in
// ~/.copilot/mcp-config.json (or $COPILOT_HOME/mcp-config.json).
//
// Merge-safe: reads existing config, adds ctx server, writes back.
// Skips if ctx server is already registered.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if file read/write fails
func ensureMCPConfig(cmd *cobra.Command) error {
	copilotHome := os.Getenv(cfgHook.EnvCopilotHome)
	if copilotHome == "" {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return homeErr
		}
		copilotHome = filepath.Join(home, cfgHook.DirCopilotHome)
	}

	target := filepath.Join(copilotHome, cfgHook.FileMCPConfigJSON)

	// Read existing config if it exists
	existing := make(map[string]interface{})
	if data, readErr := os.ReadFile(filepath.Clean(target)); readErr == nil {
		if jErr := json.Unmarshal(data, &existing); jErr != nil {
			return jErr
		}
	}

	// Get or create mcpServers map
	servers, _ := existing[cfgHook.KeyMCPServers].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	// Check if ctx is already registered
	if _, ok := servers[mcpServer.Name]; ok {
		writeSetup.InfoCopilotCLISkipped(cmd, target)
		return nil
	}

	// Add ctx MCP server
	servers[mcpServer.Name] = map[string]interface{}{
		cfgHook.KeyType:    cfgHook.MCPServerType,
		cfgHook.KeyCommand: mcpServer.Command,
		cfgHook.KeyArgs:    mcpServer.Args(),
		cfgHook.KeyTools:   []string{cfgHook.ToolsWildcard},
	}
	existing[cfgHook.KeyMCPServers] = servers

	// Create directory if needed
	if mkdirErr := io.SafeMkdirAll(copilotHome, fs.PermExec); mkdirErr != nil {
		return mkdirErr
	}

	data, marshalErr := json.MarshalIndent(existing, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	data = append(data, token.NewlineLF...)

	if writeFileErr := io.SafeWriteFile(target, data, fs.PermFile); writeFileErr != nil {
		return writeFileErr
	}
	writeSetup.InfoCopilotCLICreated(cmd, target)
	return nil
}
