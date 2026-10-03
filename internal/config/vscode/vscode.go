//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package vscode

// Dir is the VS Code workspace configuration directory.
const Dir = ".vscode"

// FileMCPJSON is the MCP server configuration file name within .vscode/.
const FileMCPJSON = "mcp.json"

// JSON keys for mcp.json.
const (
	KeyServers = "servers"
	KeyCommand = "command"
	KeyArgs    = "args"
)
