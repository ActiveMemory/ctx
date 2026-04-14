//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package handler holds the **domain logic** behind every MCP
// (Model Context Protocol) tool that ctx exposes to Claude
// Code and other MCP-compatible clients.
//
// The package is intentionally protocol-free. Every exported
// function takes typed Go parameters (a `*entity.MCPDeps`, a
// path, a string, a struct) and returns `(string, error)` —
// the formatted user-facing reply and a Go error. The sister
// package [internal/mcp/server] handles JSON-RPC framing,
// argument extraction from `map[string]any`, and response
// wrapping. This split keeps the domain logic
// **unit-testable without standing up a server** and makes it
// reusable from non-MCP callers (notably the CLI's
// `ctx agent`).
//
// # The Tool Surface
//
// The functions in [tool.go] correspond one-to-one with the
// MCP tools advertised by the server. A non-exhaustive
// inventory:
//
//   - [Status]                 — context summary (file list,
//     token counts, drift signals).
//   - **`ctx_add`**            — add a task / decision /
//     learning / convention.
//   - **`ctx_complete`**       — flip a task from `[ ]` to
//     `[x]` via [taskComplete].
//   - **`ctx_compact`**        — invoke [tidy] to archive
//     done work.
//   - **`ctx_drift`**          — run [drift.Detect] and
//     render the report.
//   - **`ctx_journal_source`** — list raw session
//     transcripts via [journal/parser].
//   - **`ctx_search`**         — text search across context
//     files via [internal/entry].
//   - **`ctx_remind`**         — read/dismiss reminders via
//     [remindStore].
//   - **`ctx_session_*`**      — `session_start`,
//     `session_end`, `session_event` lifecycle plumbing
//     (covered in [session_hooks.go]).
//   - **`ctx_steering_get`**   — surface matched steering
//     files via [steering.go] (see [internal/steering]).
//   - **`ctx_check_task_completion`** — match recent file
//     edits to open tasks.
//   - **`ctx_watch_update`**   — apply context updates the
//     agent emits in `<ctx-update>` blocks.
//
// Each function loads context fresh via [load.Do] when it
// needs current state — there is no per-tool cache. This
// keeps the response correct after edits the agent itself
// just made.
//
// # Governance — The Append-on-Every-Reply Layer
//
// [governance.go] implements the **governance trailer**:
// short, structured warnings that ride along with every MCP
// reply when the session has accumulated overdue work.
// [CheckGovernance] is invoked by the server **after** the
// tool has produced its answer; it consults the per-session
// state on `entity.MCPDeps`, drains the VS Code extension's
// violations file ([violations.go]), and assembles a
// newline-separated banner of nudges to append.
//
// The function is a free function rather than a method on
// `MCPSession` precisely because it does I/O (reading the
// violations file). `toolName` is passed in so the function
// can suppress redundant warnings — e.g. the drift warning
// is not appended to a `ctx_drift` response, since the user
// is already looking at it.
//
// # Violations Drain
//
// The Claude Code VS Code extension records hook-detected
// violations to a JSON file under the context dir. The
// handler reads it with [readViolations], surfaces the
// entries, and **truncates the file** so each violation
// surfaces exactly once. The JSON shape is
// [violationsData] / [violation].
//
// # Session Hooks
//
// [session_hooks.go] implements the three lifecycle tools
// (`session_start`, `session_end`, `session_event`) the MCP
// client calls to mark transitions. They write to per-session
// state files under `state/` and emit nudge messages when the
// configured ceremonies have been skipped.
//
// # Concurrency
//
// Handler functions are reentrant; they hold no module-level
// state. Per-session state lives on [entity.MCPDeps] (passed
// in by the server) and on the per-session files in `state/`,
// which are written through the package's own append helpers.
//
// # Related Packages
//
//   - [github.com/ActiveMemory/ctx/internal/mcp/server] —
//     JSON-RPC framing, tool dispatch, response wrapping.
//   - [github.com/ActiveMemory/ctx/internal/mcp/proto] —
//     wire-protocol types.
//   - [github.com/ActiveMemory/ctx/internal/entity] —
//     `MCPDeps`, `MCPSession`, and the domain types these
//     handlers operate on.
//   - [github.com/ActiveMemory/ctx/internal/context/load] —
//     reads and assembles the loaded [entity.Context].
//   - [github.com/ActiveMemory/ctx/internal/drift],
//     [internal/tidy], [internal/journal/parser],
//     [internal/steering] — domain packages this handler
//     orchestrates.
package handler
