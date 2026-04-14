//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements **`ctx agent`** — the command
// that produces an AI-ready, token-budgeted context
// packet for injection into the next prompt.
//
// `ctx agent` is the most-called user-facing command in
// production: tool integrations (Claude Code's
// `PreToolUse` hook, Copilot CLI's session-start hook,
// the Cursor MCP server) all invoke it on every prompt
// to assemble what the AI sees.
//
// # Public Surface
//
//   - **[Cmd]** — cobra command with `--budget N`
//     (default 8000), `--format markdown|json`,
//     `--include-hub`, and the `--prompt <text>`
//     companion that lets the budget allocator score
//     for relevance against the user's actual prompt.
//   - **[Run]** — loads context via
//     [internal/context/load], optionally folds in
//     hub entries (`--include-hub`), assembles the
//     packet via [internal/cli/agent/core/budget],
//     scores entries via
//     [internal/cli/agent/core/score], and renders
//     to stdout.
//
// # Performance
//
// The whole call typically completes in 50–150 ms on
// a project with hundreds of entries. The cost is
// dominated by file IO (the per-file token estimator
// is fast), which is why
// [internal/context/load] reads the smallest set of
// files needed and the budget allocator stops as
// soon as the budget is exhausted.
//
// # Concurrency
//
// Single-process, sequential.
//
// # Related Packages
//
//   - [internal/cli/agent/core/budget]   — the
//     allocator that decides what makes the cut.
//   - [internal/cli/agent/core/score]    — entry
//     relevance scorer.
//   - [internal/context/load]            — the
//     context-loading layer.
//   - [internal/cli/hub] /
//     [internal/cli/connection]          — supply
//     hub entries for `--include-hub`.
package root
