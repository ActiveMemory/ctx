//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ping responds to MCP ping requests with an
// empty success result.
//
// # Handler
//
// Dispatch handles the "ping" JSON-RPC method by
// returning an empty success response. MCP clients
// use ping to verify that the server is alive and
// responsive.
//
//	resp := ping.Dispatch(req)
//
// # Protocol
//
// The ping method belongs to the MCP base protocol.
// It takes no parameters and returns an empty object.
// The response echoes the request ID so the client
// can correlate it with the original request.
package ping
