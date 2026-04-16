//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package kiro generates Kiro editor integration files
// during project setup.
//
// Kiro is an AI-native IDE. This package creates the
// configuration files that connect Kiro to the ctx MCP
// server and synchronize steering rules.
//
// # Deployment Steps
//
// [Deploy] performs two operations in sequence:
//
//  1. MCP configuration: creates
//     .kiro/settings/mcp.json with the ctx MCP server
//     entry. The JSON structure uses a top-level
//     "mcpServers" key and includes additional fields
//     for disabled state and auto-approve tool lists.
//     Skips if the file already exists.
//  2. Steering sync: copies steering files from the
//     context directory to .kiro/steering/ in Kiro's
//     native markdown format.
//
// Both operations delegate to the shared mcp package
// for the actual file writing and steering sync logic.
//
// # Auto-Approve
//
// The Kiro config includes an autoApprove list of MCP
// tool names that Kiro can invoke without user
// confirmation. This includes read-only tools like
// status, steering-get, search, session-start,
// session-end, next, and remind.
//
// # Types
//
// The package defines mcpConfig and serverEntry types
// that model the .kiro/settings/mcp.json structure.
// The serverEntry type extends the base pattern with
// Disabled and AutoApprove fields.
package kiro
