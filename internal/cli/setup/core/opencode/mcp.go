//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"encoding/json"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureMCPConfig registers the ctx MCP server in opencode.json
// at the project root.
//
// Merge-safe: reads existing config, adds ctx server under
// the "mcp" key, writes back. Skips if ctx server is already
// registered.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if file read/write fails
func ensureMCPConfig(cmd *cobra.Command) error {
	target := cfgHook.FileOpenCodeJSON

	// Read existing config if it exists.
	existing := make(map[string]interface{})
	data, readErr := ctxIo.SafeReadUserFile(target)
	if readErr == nil {
		if jErr := json.Unmarshal(data, &existing); jErr != nil {
			return jErr
		}
	}

	// Get or create mcp map.
	servers, _ := existing[cfgHook.KeyMCP].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	// Check if ctx is already registered.
	if _, ok := servers[mcpServer.Name]; ok {
		writeSetup.InfoOpenCodeSkipped(cmd, target)
		return nil
	}

	// Add ctx MCP server.
	servers[mcpServer.Name] = map[string]interface{}{
		cfgHook.KeyType:    cfgHook.MCPServerType,
		cfgHook.KeyCommand: mcpServer.Command,
		cfgHook.KeyArgs:    mcpServer.Args(),
	}
	existing[cfgHook.KeyMCP] = servers

	data, marshalErr := json.MarshalIndent(
		existing, "", token.Indent2,
	)
	if marshalErr != nil {
		return marshalErr
	}
	data = append(data, token.NewlineLF...)

	writeFileErr := ctxIo.SafeWriteFile(
		target, data, fs.PermFile,
	)
	if writeFileErr != nil {
		return writeFileErr
	}
	writeSetup.InfoOpenCodeCreated(cmd, target)
	return nil
}
