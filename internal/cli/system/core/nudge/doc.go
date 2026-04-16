//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package nudge holds the **shared nudge-emission helpers**
// every `ctx system check_*` hook calls when it has decided to
// surface a message to the user (or to the agent through the
// VERBATIM relay path).
//
// The package is the *muscle*; each `cmd/check_*` package is
// the *brain*. The check decides when a nudge fires; this
// package decides what it looks like and where it lands. That
// split keeps the per-check files small and ensures that
// every nudge (checkpoint, oversize, billing window, pause
// banner) has the same shape and routing.
//
// # Emission Path
//
//   - **[EmitCheckpoint](msg)**: fires a "context
//     checkpoint reached" nudge: prompt-counter trip,
//     persistence-stale signal, etc. Routes through the
//     VERBATIM relay so the user (and the agent) both see
//     the exact text.
//   - **[EmitWindowWarning](used, total)**: fires when
//     session token usage crosses the configured
//     `injection_token_warn` (or `context_window`)
//     threshold. One-shot per session.
//   - **[EmitBillingWarning](used)**: fires the one-shot
//     "you've exceeded your included token allowance"
//     nudge for Claude Pro 1M-context users; gated by
//     `billing_token_warn` in `.ctxrc`.
//
// All three honor the **session pause** flag (see
// [Paused]) so a user who has explicitly silenced ceremony
// nudges sees nothing, except for the security-relevant
// hooks, which fire regardless.
//
// # Pause Semantics
//
// [PauseMarkerPath](sessionID) returns the per-session
// marker file. [Pause] / [Resume] write/remove it
// (exported for use by `ctx hook pause` /
// `ctx hook resume`). [Paused] returns the configured
// state and the turn count since pause began so the
// graduated reminder can render `ctx:paused (N)`.
//
// # Trigger Evaluation
//
// [EvaluateTrigger] is the per-check predicate
// dispatcher: takes a check name + threshold, reads the
// per-session counter, decides fire-or-not, increments
// the counter atomically. It is what every `cmd/check_*`
// hook calls before emitting.
//
// # Concurrency
//
// All functions are stateless or transact through the
// per-session marker files (which are atomic on POSIX
// for the small writes we issue). Concurrent hook fires
// across separate sessions never race.
package nudge
