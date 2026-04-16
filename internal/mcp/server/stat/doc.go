//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stat provides session statistics helpers
// for the MCP server.
//
// # Counting
//
// TotalAdds sums all entry-add counts from a map
// keyed by entry type. This is used during session
// checkpoint reporting to produce a single total
// across all types (tasks, decisions, learnings,
// conventions, etc.).
//
//	counts := map[string]int{
//	    "task": 3, "decision": 1,
//	}
//	total := stat.TotalAdds(counts) // => 4
//
// # Design
//
// The function accepts a generic map[string]int so
// it remains decoupled from the specific entry types
// defined elsewhere. Callers maintain their own count
// maps and pass them in when needed.
package stat
