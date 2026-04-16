//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp provides the "ctx mcp" CLI command for
// starting the Model Context Protocol server.
//
// The MCP server exposes ctx context operations as MCP
// tools that AI coding assistants can invoke over stdio
// transport. This allows tools like Claude Code, Cursor,
// and other MCP-aware clients to read, write, and query
// project context without shelling out to the ctx CLI.
//
// # Subpackages
//
//	cmd/root -- MCP server bootstrap, tool registration,
//	  and stdio transport setup. Starts the MCP server on
//	  stdio, registering tool handlers for context
//	  operations. The command annotates itself with SkipInit
//	  so it can run without a fully initialized .context/
//	  directory.
package mcp
