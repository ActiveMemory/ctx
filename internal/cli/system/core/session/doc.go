//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session provides the **shared session-state
// helpers** every `ctx system` hook calls when it needs to
// (a) read the JSON envelope Claude Code sent on stdin,
// (b) extract the session ID, or (c) write the session-stats
// counters that drive ceremony nudges.
//
// The package owns the per-session state files under
// `.context/state/session-<id>.json`, the lightweight
// counters that hooks like `check_ceremony`,
// `check_persistence`, and `check_context_size` evaluate
// each time they fire.
//
// # Public Surface
//
//   - **[ReadInput]**: decodes the JSON envelope Claude
//     Code writes to a hook's stdin. Returns a typed
//     [HookInput] regardless of which event fired.
//   - **[ReadID]**: convenience; pulls just the
//     session ID from the envelope.
//   - **[FormatContext](payload)**: formats a JSON
//     payload as the canonical "context block" that
//     hooks emit on stdout to inject content into the
//     agent's next prompt.
//   - **[LatestPct](contextDir)**: returns the most
//     recent context-window-usage percentage written by
//     the agent CLI. Used by the size-checkpoint hook
//     to decide whether to nudge.
//   - **[WriteStats](sessionID, stats)**: atomically
//     updates the per-session stats file. Each `check_*`
//     hook owns a slice of the stats struct.
//
// # State File Lifecycle
//
// `session-<id>.json` is created lazily on first hook
// fire, updated in place by subsequent hooks, and pruned
// by `ctx prune` after the session has been idle for
// the configured threshold.
//
// # Concurrency
//
// Hooks within a session fire serially (Claude Code
// drives them one at a time per turn) so the in-process
// concurrency model is single-writer. The atomic
// rename pattern in [WriteStats] guards the rare case
// where a sibling process inspects the file mid-write.
package session
