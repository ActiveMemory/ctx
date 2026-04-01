//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resource handles MCP resource requests including list,
// read, subscribe, and unsubscribe operations.
//
// Each dispatcher validates parameters, delegates to the
// appropriate handler, and returns JSON-RPC 2.0 responses.
//
// Key exports: [DispatchList], [DispatchRead],
// [DispatchSubscribe], [DispatchUnsubscribe].
// Part of the MCP server (JSON-RPC 2.0 over stdin/stdout).
package resource
