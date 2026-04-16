//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package detect resolves the **reference timestamp**
// `ctx change` uses to compute "what changed since". The
// CLI offers three signals; this package picks the best
// one and returns it.
//
// `ctx change` answers "what moved since I was last in
// this project?": context file edits, code commits,
// directories touched. Picking the right "since when"
// is the package's only job.
//
// # The Three Signals
//
// In priority order:
//
//  1. **Explicit `--since`**: `ctx change --since
//     2026-04-12` or `--since 3d`. Parsed by
//     [ParseSinceFlag] into a time.Time.
//  2. **Session marker**: `[FromMarkers]` reads
//     `state/session-event.jsonl` for the timestamp
//     of the last session-end event. The most useful
//     "since" for "since I was last here".
//  3. **Event log**: `[FromEvents]` falls back to
//     the newest hook event timestamp when no
//     session-end marker exists.
//
// [ReferenceTime] composes the three: returns the
// `--since` value when set; otherwise the more recent
// of the marker / event timestamps; falls back to
// "30 days ago" when nothing is known.
//
// # Flag Parsing
//
// [ParseSinceFlag] accepts:
//
//   - **Date**: `2026-04-12` (parsed midnight
//     UTC).
//   - **Duration**: `3d`, `12h`, `2w`, `1m` (rich
//     duration syntax beyond Go's stdlib).
//   - **`yesterday`**, **`today`**: relative
//     keywords.
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent
// callers never race.
package detect
