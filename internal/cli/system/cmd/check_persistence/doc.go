//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_persistence implements the
// **`ctx system check-persistence`** hook, the nudge
// that tells the user (and the agent) "you have done a
// lot of work without persisting anything to
// `.context/`; consider adding a learning, decision, or
// task update before the session ends".
//
// The hook tracks **prompts since the last `.context/`
// file mtime change**: every time the user submits a
// prompt, the hook increments a counter; every time a
// `.context/` file is touched, the counter resets.
// When the counter crosses the configured threshold,
// the hook emits the persistence reminder via the
// VERBATIM relay.
//
// # Public Surface
//
//   - **[Cmd]**: cobra command (hidden under
//     `ctx system`).
//   - **[Run]**: reads the JSON envelope, scans
//     `.context/` for the most recent mtime,
//     compares against the prompt counter, and
//     fires the nudge when the threshold is
//     crossed.
//
// # Why "Mtime" Not "Edits"
//
// Mtime is the simplest proxy for "the user
// captured something" that does not require
// instrumenting every write path. It catches
// `ctx add` writes, the agent's direct edits,
// and even out-of-band `vi` edits.
//
// # Concurrency
//
// Single-process per session.
package check_persistence
