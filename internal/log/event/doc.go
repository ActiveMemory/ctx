//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package event implements the **JSONL hook event log** —
// the append-only on-disk record of every hook lifecycle
// event ctx generates so users can inspect, audit, or
// timeline what happened in a session.
//
// Two pieces of the system depend on it:
//
//   - **`ctx hook event`** — user-facing query: "what did
//     the hooks do during the last session?".
//   - **`ctx system check_persistence`** and friends — read
//     the log to detect "you committed but never wrote a
//     decision" patterns and nudge accordingly.
//
// # On-Disk Format
//
// The log lives at `.context/state/events.jsonl` and is
// **append-only JSONL**: one [Event] per line, written via
// [Append], rotated to `events.1.jsonl` when the file
// exceeds [config/event.LogMaxBytes] (1 MiB). At most one
// rotation generation is kept; older history is discarded.
//
// # Opt-In
//
// Logging is **disabled by default** — many users do not
// want hook activity persisted. [Append] is a noop when
// `event_log: false` in `.ctxrc`; setting it to `true`
// activates collection. The `ctx hook event` query
// gracefully reports "no events recorded" when the file is
// missing.
//
// # The Query Surface
//
// [Query](opts) reads both `events.jsonl` and the rotated
// `events.1.jsonl` (in chronological order), then applies
// the filters from [entity.EventQueryOpts]:
//
//   - **Hook** — match a specific hook name
//     (e.g. `check-persistence`).
//   - **Session** — match a session ID prefix.
//   - **Event** — match an event-type tag (`fired`,
//     `relayed`, `blocked`, …).
//   - **Last N** — keep only the most recent N matches
//     (default [config/event.DefaultLast] = 50).
//
// # Concurrency
//
// [Append] uses an O_APPEND open which is atomic for
// small (sub-PIPE_BUF) writes on POSIX systems — the
// log line size we emit is well under that bound, so
// concurrent appenders interleave but never tear a line.
// [Query] reads a snapshot of the file; concurrent
// appends mid-read are tolerated (the worst case is a
// half-written final line that the JSONL decoder skips).
package event
