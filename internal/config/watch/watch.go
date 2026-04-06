//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

// Stream scanner buffer sizes.
const (
	// StreamScannerInitCap is the initial capacity for the scanner buffer.
	StreamScannerInitCap = 64 * 1024
	// StreamScannerMaxSize is the maximum size for the scanner buffer.
	StreamScannerMaxSize = 1024 * 1024
)

// XML attribute extraction constants.
const (
	// ContextUpdateMinGroups is the minimum number of regex capture
	// groups expected from a context-update match (full match + tag + content).
	ContextUpdateMinGroups = 3
)

// Default provenance for watch-originated entries.
//
// Watch streams are machine-generated, so entries receive
// fixed provenance identifying the source as the watch pipeline.
const (
	// ProvenanceSessionID is the default session ID for watch entries.
	ProvenanceSessionID = "watch"
	// ProvenanceBranch is the default branch for watch entries.
	ProvenanceBranch = "watch"
	// ProvenanceCommit is the default commit for watch entries.
	ProvenanceCommit = "watch"
)
