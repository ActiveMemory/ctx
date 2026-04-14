//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package proto defines the **JSON-RPC 2.0** wire types and
// the **MCP** (Model Context Protocol) extension types ctx
// speaks with Claude Code and other MCP-compatible clients.
//
// The package is wire-protocol-only: structs with `json:`
// tags, sentinel constants for method names and error codes,
// no logic. Domain behavior lives in
// [internal/mcp/handler]; transport in [internal/mcp/server].
//
// # JSON-RPC 2.0
//
//   - **[Request]** — `jsonrpc`, `id`, `method`,
//     `params`. `id` may be string, number, or null per
//     the spec.
//   - **[Response]** — `jsonrpc`, `id`, plus exactly
//     one of `result` / `error`.
//   - **[Error]** — `code`, `message`, optional `data`.
//     The standard error codes ([CodeParseError],
//     [CodeInvalidRequest], etc.) are exported as
//     constants.
//   - **[Notification]** — `Request` without an `id`,
//     used for one-way messages (logging, progress).
//
// # MCP Extensions
//
// MCP layers these methods on top of JSON-RPC:
//
//   - **`tools/list` / `tools/call`** — for tool
//     dispatch.
//   - **`prompts/list` / `prompts/get`** — for
//     server-curated prompts.
//   - **`resources/list` / `resources/read` /
//     `resources/subscribe`** — for server-exposed
//     resources.
//
// Each method has a typed request and response struct
// in this package: [ToolsCallRequest],
// [ToolsCallResponse], [Tool], [PromptsGetResponse],
// etc.
//
// # Stability
//
// The wire shape is fixed by external specifications
// (JSON-RPC 2.0 and the MCP spec). Changes here
// require coordinated client updates and should not
// happen casually. The audit suite watches for
// accidental field renames.
//
// # Concurrency
//
// All exports are immutable types. Encoding /
// decoding is goroutine-safe at the
// `encoding/json` boundary.
//
// # Related Packages
//
//   - [internal/mcp/server]   — encodes and decodes
//     these types; calls into [internal/mcp/handler]
//     for the result.
//   - [internal/mcp/handler]  — produces the typed
//     payloads the server marshals into [Response].
package proto
