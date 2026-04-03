//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package copilot deploys GitHub Copilot integration files.
//
// [DeployInstructions] generates .github/copilot-instructions.md and
// the accompanying .vscode/mcp.json for VS Code Copilot MCP support.
// These files give Copilot access to project context through the ctx
// MCP server, enabling context-aware completions and chat responses.
//
// Key exports: [DeployInstructions].
// Called by the setup core orchestrator during ctx init.
package copilot
