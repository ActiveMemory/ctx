//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package summary generates the **brief, human-readable
// summaries** ctx uses in compact-output contexts: the
// `ctx status` one-line health line, the agent context
// packet's lead paragraph, the doctor's roll-up banner.
//
// The summaries operate against an already-loaded
// [entity.Context] so they share data with every other
// downstream consumer, avoiding double-reading the filesystem.
//
// # Public Surface
//
//   - **[Generate](ctx, opts)**: produces a multi-line
//     summary string covering:
//   - file count + total token estimate
//   - newest / oldest file mtime
//   - drift signal counts (warnings, violations)
//   - open task count
//   - last session timestamp from
//     `state/session-event.jsonl`
//     Output shape is tunable via [opts] for the
//     different consumers (`ctx status` wants a single
//     line; the agent packet wants 3-5 lines).
//
// # Why a Dedicated Package
//
// Three callers need the same numbers (`ctx status`,
// `ctx agent`, `ctx doctor`) and three different
// renderings. Hoisting the data computation here means
// each renderer reuses the same byte / token / mtime
// counters, and a fix to the count logic only has to
// happen once.
//
// # Concurrency
//
// All functions are pure data transformations over
// the input [entity.Context]. Concurrent callers
// never race.
package summary
