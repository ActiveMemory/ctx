//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package opencode generates OpenCode integration files during
// project setup.
//
// OpenCode is a terminal-first AI coding agent (opencode.ai).
// This package creates the configuration files that connect
// OpenCode to the ctx MCP server, deploy a thin lifecycle
// plugin, and synchronize skills.
//
// # Deployment Steps
//
// [Deploy] performs four operations in sequence:
//  1. Plugin deployment: creates .opencode/plugins/ctx/ with
//     index.ts and package.json
//  2. MCP configuration: merges ctx server into opencode.json
//  3. AGENTS.md: deploys shared agent instructions
//  4. Skills: copies ctx skills to .opencode/skills/
package opencode
