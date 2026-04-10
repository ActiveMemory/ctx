//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// MCPDeps bundles the ambient runtime inputs that every MCP handler
// function needs. It is held once by the MCP server and threaded
// through dispatch into each tool/prompt implementation.
//
// Handler functions in internal/mcp/handler accept *MCPDeps as their
// first argument. The struct replaces an earlier god-object Handler
// type whose only real job was to carry these three fields around.
//
// Fields:
//   - ContextDir: Absolute path to the .context/ directory
//   - TokenBudget: Maximum token budget for context assembly
//   - Session: Per-run advisory state (governance counters, pending
//     updates, etc.)
type MCPDeps struct {
	ContextDir  string
	TokenBudget int
	Session     *MCPSession
}
