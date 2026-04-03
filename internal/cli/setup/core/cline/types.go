//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cline

// vscodeMCPConfig is the top-level mcp.json structure for Cline.
type vscodeMCPConfig struct {
	Servers map[string]vscodeMCPServer `json:"servers"`
}

// vscodeMCPServer describes one MCP server entry in mcp.json.
type vscodeMCPServer struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}
