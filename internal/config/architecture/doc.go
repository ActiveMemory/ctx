//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package architecture provides constants that drive the
// architecture mapping and staleness detection subsystem.
//
// ctx maintains an architecture map (ARCHITECTURE.md) that
// documents module boundaries, dependencies, and ownership.
// This package defines the file names, staleness thresholds,
// and template variables used by the mapping pipeline.
//
// # Map Tracking
//
// [MapTracking] names the JSON state file that records
// which modules have been covered by the architecture
// mapper. The file lives in .context/ and is updated each
// time the /ctx-architecture skill completes a pass.
//
// # Staleness Detection
//
// Architecture maps drift as code evolves. The staleness
// hook nudges agents to refresh when the map is outdated:
//
//   - [MapStaleDays] sets the threshold to 30 days. If the
//     last refresh is older, a nudge fires.
//   - [MapStalenessThrottleID] names the daily throttle
//     state file so the nudge fires at most once per day.
//
// # Template Variables
//
// Hook messages are rendered from Go templates. The
// architecture hooks inject these variables:
//
//   - [VarLastRefreshDate]: the ISO date of the last
//     architecture refresh.
//   - [VarModuleCount]: the number of modules that have
//     changed since the last refresh.
//
// # Why Centralized
//
// The staleness threshold, throttle ID, and template keys
// are consumed by both the hook runner and the skill
// implementation. Keeping them here avoids import cycles
// and makes the values discoverable.
package architecture
