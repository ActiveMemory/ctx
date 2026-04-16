//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize handles the MCP initialize
// handshake that establishes the server session.
//
// # Handler
//
// Dispatch responds to the "initialize" JSON-RPC
// method by returning the server's capabilities and
// version information. The response includes:
//
//   - ProtocolVersion: the MCP protocol version
//     supported by this server.
//   - Capabilities: resource subscriptions, tools,
//     and prompts support flags.
//   - ServerInfo: the server name and version string.
//
// # Usage
//
//	resp := initialize.Dispatch(version, req)
//
// # Protocol
//
// The initialize method is the first request an MCP
// client sends. It must complete before the client
// can call any other methods. The server advertises
// which capabilities it supports so the client knows
// what methods are available.
package initialize
