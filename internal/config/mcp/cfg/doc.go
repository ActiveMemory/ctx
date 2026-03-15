//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cfg defines MCP server defaults and limits.
//
// These constants are used exclusively by the MCP server
// (internal/mcp) and are not configurable via .ctxrc. The
// truncation lengths (TruncateLen, TruncateContentLen) mirror
// values used in compact CLI output; if those become
// configurable, they should be unified into a shared location.
package cfg
