//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cursor generates Cursor editor integration
// files during project setup.
//
// Cursor is an AI-native code editor. This package
// creates the configuration files that connect Cursor
// to the ctx MCP server and synchronize steering rules.
//
// # Deployment Steps
//
// [Deploy] performs two operations in sequence:
//
//  1. MCP configuration: creates .cursor/mcp.json with
//     the ctx MCP server entry. The JSON structure uses
//     a top-level "mcpServers" key containing server
//     name, command, and arguments. Skips if the file
//     already exists.
//  2. Steering sync: copies steering files from the
//     context directory to .cursor/rules/ in Cursor's
//     native MDC format.
//
// Both operations delegate to the shared mcp package
// for the actual file writing and steering sync logic.
//
// # Types
//
// The package defines mcpConfig and serverEntry types
// that model the .cursor/mcp.json structure. The key
// difference from Cline is the "mcpServers" JSON key
// instead of "servers".
//
// # Data Flow
//
// The setup orchestrator calls [Deploy]. It builds an
// mcpConfig struct with server details from
// config/mcp/server, then passes it to mcp.Deploy for
// file creation. Steering sync is handled by
// mcp.SyncSteering with the Cursor tool identifier.
package cursor
