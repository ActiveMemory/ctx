//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_map_staleness implements the
// **`ctx system check-map-staleness`** hidden hook,
// which detects when the architecture map has not been
// refreshed for too long.
//
// # What It Does
//
// The hook reads the map-tracking.json file to find
// when the architecture map was last regenerated. If
// the last-run date exceeds the configured staleness
// window (typically 14 days), it counts how many
// commits have touched internal/ since then. When both
// conditions are met (stale age + relevant commits),
// it emits a nudge recommending that the user rerun
// the architecture mapping command.
//
// If map tracking is opted out or no relevant commits
// have been made, the hook stays silent.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On stale map with relevant commits: a nudge showing
// the last-run date and commit count. On fresh map, no
// commits, opted out, or throttled: no output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], loads the tracking
// data through [core/health.ReadMapTracking], counts
// commits via [core/health.CountModuleCommits], and
// emits the warning through
// [core/health.EmitMapStalenessWarning].
package check_map_staleness
