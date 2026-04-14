//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cfg

// MCP server configuration constants.
const (
	// ScanMaxSize is the maximum scanner buffer size for MCP messages (1 MB).
	ScanMaxSize = 1 << 20

	// DefaultSourceLimit is the max sessions returned by ctx_journal_source.
	DefaultSourceLimit = 5
	// MaxSourceLimit caps the source limit to prevent unbounded queries.
	MaxSourceLimit = 100
	// MinWordLen is the shortest word considered for overlap matching.
	MinWordLen = 4
	// MinWordOverlap is the minimum word matches to signal task completion.
	MinWordOverlap = 2

	// --- Input length limits (MCP-SAN.1) ---

	// MaxContentLen is the maximum byte length for entry content fields.
	MaxContentLen = 32_000
	// MaxNameLen is the maximum byte length for tool/prompt/resource names.
	MaxNameLen = 256
	// MaxQueryLen is the maximum byte length for search queries.
	MaxQueryLen = 1_000
	// MaxCallerLen is the maximum byte length for caller identifiers.
	MaxCallerLen = 128
	// MaxURILen is the maximum byte length for resource URIs.
	MaxURILen = 512
)
