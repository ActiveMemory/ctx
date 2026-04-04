//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kiro

// mcpConfig is the JSON structure for .kiro/settings/mcp.json.
type mcpConfig struct {
	MCPServers map[string]serverEntry `json:"mcpServers"`
}

// serverEntry describes one MCP server entry in mcp.json.
type serverEntry struct {
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Disabled    bool     `json:"disabled"`
	AutoApprove []string `json:"autoApprove"`
}
