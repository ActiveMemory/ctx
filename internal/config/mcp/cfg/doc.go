//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cfg defines operational constants that tune the
// MCP server's runtime behavior -- buffer sizes, default
// query limits, and word-matching thresholds.
//
// These values are compile-time constants, not user-facing
// configuration. They are consumed exclusively by the MCP
// server (internal/mcp) and its handlers.
//
// # Key Constants
//
//   - [ScanMaxSize] (1 MB) -- the maximum buffer
//     allocated by the JSON-RPC scanner when reading
//     a single MCP message from stdin. Messages
//     larger than this are rejected.
//   - [DefaultSourceLimit] (5) -- the default cap on
//     sessions returned by the ctx_journal_source
//     tool when the caller omits a limit.
//   - [MinWordLen] (4) -- the shortest word considered
//     when computing overlap between a recent action
//     description and a pending task title.
//   - [MinWordOverlap] (2) -- the minimum number of
//     matching words required to signal that a task
//     has likely been completed by the recent action.
//
// # Why These Are Centralized
//
// Handlers, the scanner, and the task-completion
// heuristic each reference these values. Keeping them
// in a single file prevents silent divergence when one
// handler is tuned but another is not.
package cfg
