//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cfg

const (
	// ScanMaxSize is the maximum scanner buffer size for MCP messages (1 MB).
	ScanMaxSize = 1 << 20

	// DefaultSourceLimit is the max sessions returned by ctx_journal_source.
	DefaultSourceLimit = 5
	// MinWordLen is the shortest word considered for overlap matching.
	MinWordLen = 4
	// MinWordOverlap is the minimum word matches to signal task completion.
	MinWordOverlap = 2
)
