//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx mcp" command.
//
// # Overview
//
// The mcp command starts a Model Context Protocol (MCP)
// server that communicates over stdin/stdout using
// JSON-RPC 2.0. It exposes ctx's context management
// capabilities as MCP tools, allowing AI assistants to
// read and manipulate project context programmatically.
//
// This command is typically invoked by an MCP client
// (such as Claude Code) rather than run directly by the
// user. The client launches "ctx mcp" as a subprocess
// and communicates via the standard streams.
//
// # Flags
//
// This command accepts no flags. It reads the context
// directory from rc and the version from the root
// cobra.Command.
//
// # Behavior
//
// [Cmd] creates a new MCP server instance using the
// resolved context directory and the CLI version string,
// then calls srv.Serve which blocks until the client
// disconnects or an I/O error occurs.
//
// The server registers tools for reading context files,
// querying project state, and other context operations
// defined in the internal/mcp/server package.
//
// # Output
//
// All communication happens over stdin/stdout in
// JSON-RPC 2.0 format. No human-readable output is
// produced on stderr under normal operation.
package root
