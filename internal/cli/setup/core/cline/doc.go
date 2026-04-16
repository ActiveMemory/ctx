//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cline generates Cline editor integration files
// during project setup.
//
// Cline is a VS Code extension for AI-assisted coding.
// This package creates the configuration files that
// connect Cline to the ctx MCP server and synchronize
// steering rules.
//
// # Deployment Steps
//
// [Deploy] performs two operations in sequence:
//
//  1. MCP configuration: creates .vscode/mcp.json with
//     the ctx MCP server entry. The JSON structure uses
//     a top-level "servers" key containing server name,
//     command, and arguments. Skips if the file exists.
//  2. Steering sync: copies steering files from the
//     context directory to .clinerules/ in Cline's
//     native markdown format.
//
// Both operations delegate to the shared mcp package
// for the actual file writing and steering sync logic.
//
// # Types
//
// The package defines vscodeMCPConfig and
// vscodeMCPServer types that model the .vscode/mcp.json
// structure specific to Cline's format.
//
// # Data Flow
//
// The setup orchestrator calls [Deploy]. It builds a
// vscodeMCPConfig struct with server details from
// config/mcp/server, then passes it to mcp.Deploy for
// file creation. Steering sync is handled by
// mcp.SyncSteering with the Cline tool identifier.
package cline
