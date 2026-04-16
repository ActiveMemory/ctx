//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package server implements the **Model Context Protocol
// (MCP) server** that exposes ctx context, commands, and
// session-lifecycle hooks to MCP-compatible AI clients,
// primarily Claude Code, but also any other tool that speaks
// the same JSON-RPC 2.0 dialect.
//
// The server runs over **stdin/stdout** as a sub-process
// launched by the AI client; it does not bind a network port.
// Spawn behavior is configured by the client's MCP block (see
// [internal/cli/setup] for what `ctx setup` writes into each
// tool's config).
//
// # Wire Protocol
//
// MCP is JSON-RPC 2.0 with three core verbs ctx implements:
//
//   - **`tools/list`**: advertise the catalog of MCP tools
//     this server provides (`ctx_status`, `ctx_add`,
//     `ctx_complete`, `ctx_drift`, `ctx_journal_source`,
//     `ctx_search`, `ctx_steering_get`, `ctx_remind`,
//     `ctx_session_*`, `ctx_check_task_completion`,
//     `ctx_watch_update`).
//   - **`tools/call`**: invoke one tool with a typed
//     arguments map.
//   - **`prompts/list` / `prompts/get`**: surface
//     ctx-curated prompts (e.g. the session-start
//     ceremony prompt) as first-class MCP prompts.
//
// Wire types live in [internal/mcp/proto]; this package
// concerns itself with **dispatch** and **state**.
//
// # Architecture
//
// The package layers four sub-concerns:
//
//   - **[New]**: constructs a server bound to a
//     [entity.MCPDeps] (paths, runtime config). The
//     server is single-threaded by design; Claude Code
//     spawns one sub-process per session and does not
//     pipeline requests.
//   - **Routing**: [route/tool], [route/prompt],
//     [route/resource] register handlers per MCP verb.
//   - **Dispatch**: [dispatch/poll] reads one
//     JSON-RPC message at a time from stdin and routes
//     it to the right handler.
//   - **Catalog** ([catalog/data.go]): the static
//     tool/prompt/resource definitions surfaced via
//     `*/list` calls.
//
// All actual domain logic (what `ctx_drift` *does*, what
// `ctx_search` returns) lives in [internal/mcp/handler].
// This package is the protocol-aware shell around it.
//
// # Per-Session State
//
// Each running server instance owns one
// [entity.MCPSession] (turn counter, last-loaded context
// snapshot, governance flags). The session is created on
// the first `tools/call` and persists for the lifetime of
// the sub-process. The handler layer reads/mutates it
// through [entity.MCPDeps].
//
// # Governance Trailers
//
// After every `tools/call`, the dispatcher invokes
// [internal/mcp/handler.CheckGovernance] to append any
// session-overdue nudges (drift, persistence, journal
// import) to the response. The trailers are appended
// inside the JSON-RPC `result` envelope so they reach the
// AI without changing the protocol shape.
//
// # Concurrency
//
// One goroutine reads from stdin; one goroutine writes
// to stdout; tool dispatch runs in the read goroutine
// to preserve request ordering. Long-running tools
// (currently none) would need to spawn a goroutine and
// signal completion through a channel.
package server
