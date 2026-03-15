//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cfg

const (
	// ScanMaxSize is the maximum scanner buffer size for MCP messages (1 MB).
	ScanMaxSize = 1 << 20

	// MCP default values.

	DefaultRecallLimit = 5
	MinWordLen         = 4
	MinWordOverlap     = 2
	TruncateLen        = 50
	TruncateContentLen = 60
)
