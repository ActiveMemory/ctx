//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_context_size implements the
// **`ctx system check-context-size`** hook — the prompt-
// counter that fires the periodic "context checkpoint"
// nudge so users remember to wrap up their session
// before the context window fills.
//
// The hook reads the per-session prompt counter from
// `.context/state/`, increments it, and emits the
// VERBATIM checkpoint banner when the counter crosses
// any of the configured graduated thresholds (e.g.
// every 20 prompts, or at 50% / 75% / 90% of the
// configured `context_window` budget).
//
// # Public Surface
//
//   - **[Cmd]** — cobra command (hidden under
//     `ctx system`; users do not invoke this
//     directly).
//   - **[Run]** — reads the JSON envelope from
//     stdin (session ID, current usage), decides
//     whether to fire, increments the counter,
//     and writes the nudge through
//     [internal/cli/system/core/nudge.EmitCheckpoint].
//
// # Throttling
//
// To avoid nudging on every prompt, the hook
// honors the per-check throttle in
// [internal/config/hook] — at most one fire per
// configured prompt-count interval, with a
// graduated cadence as the budget pressure
// grows.
//
// # Concurrency
//
// Single-process per session. The hook is
// invoked by the AI tool's hook runtime
// serially per turn.
//
// # Related Packages
//
//   - [internal/cli/system/core/nudge]   — the
//     emission helper.
//   - [internal/cli/system/core/session] — reads
//     and writes the per-session counters.
//   - [internal/rc]                      — supplies
//     the configured `context_window` budget.
package check_context_size
