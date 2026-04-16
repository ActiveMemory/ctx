//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package out builds JSON-RPC 2.0 response structs
// for the MCP server.
//
// # Response Constructors
//
// OkResponse creates a success response with the
// given result payload. ErrResponse creates an error
// response with a code and message.
//
//	resp := out.OkResponse(id, result)
//	resp := out.ErrResponse(id, code, msg)
//
// # Tool Result Helpers
//
// ToolOK builds a successful tool result containing
// text content. ToolError builds a tool result with
// the IsError flag set. ToolResult dispatches between
// the two based on whether an error is present.
//
//	resp := out.ToolOK(id, "done")
//	resp := out.ToolError(id, "failed")
//	resp := out.ToolResult(id, text, err)
//
// # Handler Wrapper
//
// Call invokes a no-argument handler function and
// wraps its (string, error) return into a response.
// This eliminates boilerplate in tool dispatchers.
//
//	resp := out.Call(id, func() (string, error) {
//	    return "ok", nil
//	})
//
// # Content Types
//
// All tool results use the "text" content type from
// the mime config package. The JSON-RPC version string
// comes from the server config package.
package out
