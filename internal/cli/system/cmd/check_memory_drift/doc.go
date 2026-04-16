//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_memory_drift implements the
// **`ctx system check-memory-drift`** hidden hook,
// which detects drift between the agent's working
// memory and the persisted MEMORY.md file.
//
// # What It Does
//
// The hook discovers the auto-memory source path
// (typically MEMORY.md in the project root) and
// compares it against the version stored in .context/.
// When the two diverge, the hook emits a nudge
// encouraging the agent to synchronize by running the
// memory consolidation workflow.
//
// If auto-memory is not active (no source path found),
// the hook is a silent no-op.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On drift detected: a nudge message about memory
// desynchronization. On no drift or no auto-memory:
// no output.
//
// # Throttling
//
// The hook is throttled per session using a tombstone
// file (one nudge per session ID). This is different
// from the daily throttle used by most other hooks
// because memory drift is session-specific.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], discovers the
// memory source via [memory.DiscoverPath], checks for
// drift with [memory.HasDrift], and emits the nudge
// through [core/nudge.LoadAndEmit].
package check_memory_drift
