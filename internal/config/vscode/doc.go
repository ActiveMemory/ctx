//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode defines constants for generating
// VS Code workspace configuration files.
//
// `ctx setup copilot` and the cline deployer write
// .vscode/mcp.json to configure MCP server
// integration. This package provides the directory
// path, file name, and JSON keys they need.
//
// # Directory and File Paths
//
//   - [Dir] (".vscode"): the workspace config dir.
//   - [FileMCPJSON]: the MCP config file name
//     within it.
//
// # JSON Keys
//
//   - [KeyServers], [KeyCommand], [KeyArgs]:
//     mcp.json keys.
//
// # Concurrency
//
// All exports are immutable constants.
package vscode
