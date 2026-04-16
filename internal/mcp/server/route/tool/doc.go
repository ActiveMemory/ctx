//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tool is the **MCP `tools/*` dispatcher** — the
// layer that takes a JSON-RPC `tools/list` or
// `tools/call` request, validates parameters against
// the per-tool schema in [internal/mcp/server/def/tool],
// and routes to the right handler in
// [internal/mcp/handler].
//
// # Public Surface
//
//   - **[DispatchList](req, deps)** — returns the
//     full tool catalog from
//     [internal/mcp/server/def/tool.Defs].
//   - **[DispatchCall](req, deps)** — extracts the
//     tool name and arguments map from the JSON-RPC
//     params, dispatches to the matching handler,
//     wraps the handler's `(string, error)` return
//     into the MCP response envelope, then runs
//     [handler.CheckGovernance] to append any
//     overdue-work nudges.
//
// # Argument Extraction
//
// MCP tool arguments arrive as `map[string]any`
// (raw JSON). This package owns the typed
// extraction (`mustString`, `optionalInt`, etc.)
// so the handlers see typed Go values, not
// `any`.
//
// # Error Mapping
//
// Handler errors map to JSON-RPC error codes:
//
//   - **InvalidParams** — typed validation errors.
//   - **InternalError** — anything else.
//
// The original error message is included in the
// `data` field so the client can surface it to
// the user / agent.
//
// # Concurrency
//
// Sequential per server instance; see
// [internal/mcp/server].
package tool
