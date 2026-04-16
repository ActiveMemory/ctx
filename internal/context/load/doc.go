//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package load reads the `.context/` directory and
// **assembles its files into the in-memory
// [entity.Context]** that every downstream consumer —
// `ctx agent`, `ctx drift`, `ctx doctor`, every MCP
// handler — operates on.
//
// The package is the single point of truth for "what
// does the user's context look like right now?". Two
// callers running [Do] back-to-back see identical
// snapshots because the package holds no cache; freshness
// matters more than micro-optimization here.
//
// # Public Surface
//
//   - **[Do](contextDir)** — reads every file in the
//     priority order [internal/rc.PriorityOrder]
//     defines, populates an [entity.Context] with
//     each file's name, body, byte/token counts, and
//     mtime, and returns the assembled snapshot.
//
// # Read Order
//
// Files are loaded in `priority_order` from `.ctxrc`
// (default: TASKS → DECISIONS → CONVENTIONS →
// LEARNINGS → ARCHITECTURE → CONSTITUTION →
// GLOSSARY). Order matters because downstream consumers
// (notably [internal/cli/agent/core/budget]) allocate
// budget tier-by-tier in this order.
//
// # Token Counts
//
// Each file's token count is computed by the rough
// estimator in [internal/cli/agent/core/budget] —
// approximate but stable. The count is what
// [entity.Context.TokenInfo] surfaces.
//
// # Errors
//
// File-not-found for an *expected* file is silently
// tolerated (returns an empty body); the user may
// legitimately not have populated every foundation
// file yet. Read errors that are not "not found" are
// returned to the caller.
//
// # Concurrency
//
// Stateless and filesystem-bound. Concurrent
// invocations against the same directory each pay
// the full read cost — by design, see above.
package load
