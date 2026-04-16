//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package server defines the MCP server's identity
// constants, URI prefix, and launch arguments.
//
// When the ctx MCP server starts, it reports its name,
// protocol version, and capabilities to the connecting
// client during the initialize handshake. The constants
// here supply those values. The package also provides the
// [Args] helper that returns the CLI argument slice
// needed to launch the server as a subprocess.
//
// # Key Constants
//
//   - [Name] ("ctx"): the server name reported in
//     the initialize response.
//   - [Command] ("ctx"): the binary name used to
//     spawn the MCP server process.
//   - [SubcommandServe] ("serve"): the subcommand
//     under "ctx mcp" that starts the stdio server.
//   - [ResourceURIPrefix] ("ctx://context/"): the
//     URI scheme and path prefix prepended to
//     resource names to form full resource URIs.
//   - [JSONRPCVersion] ("2.0"): the JSON-RPC
//     version string included in every response.
//   - [PollIntervalSec] (5): the default interval
//     in seconds between resource-change polls.
//
// # Args Helper
//
// [Args] returns []string{"mcp", "serve"}, the argument
// slice that, when appended to the ctx binary path,
// launches the MCP server. This is consumed by the
// connect command when writing MCP client configuration
// files (e.g., .claude/settings.json).
//
// # Why These Are Centralized
//
// The initialize handler, the connect command, resource
// registration, and the polling loop all reference these
// values. A constant prevents the server from reporting
// one name while the connect command writes another.
package server
