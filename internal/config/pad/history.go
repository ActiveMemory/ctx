//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import "time"

// History directory and snapshot filename constants.
const (
	// HistoryDirName is the per-project snapshot directory, a
	// sibling of the live scratchpad file inside `.context/`.
	HistoryDirName = "scratchpad.history"

	// HistoryTimeFormat is the timestamp prefix used in
	// snapshot filenames. Sortable lexically; nanosecond
	// suffix prevents collisions across rapid-fire mutations.
	HistoryTimeFormat = "20060102T150405.000000000Z"

	// HistorySnapshotSeparator joins the timestamp and the
	// op label in a snapshot filename: `<timestamp>-<op>.enc`.
	HistorySnapshotSeparator = "-"
)

// Retention defaults for the snapshot ring buffer. Both caps
// apply: oldest snapshots are pruned once either is exceeded.
// Zero on either cap disables that cap independently.
const (
	// HistoryMaxSnapshots caps the number of retained
	// snapshots; the oldest are pruned beyond this count.
	HistoryMaxSnapshots = 20

	// HistoryMaxAge caps the age of retained snapshots;
	// snapshots older than this are pruned.
	HistoryMaxAge = 30 * 24 * time.Hour
)

// HistoryOpUnknown is the op label used when a snapshot is
// taken with no identifiable cobra subcommand context.
const HistoryOpUnknown = "write"
