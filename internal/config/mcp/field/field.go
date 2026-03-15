//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package field

// MCP tool input field names (JSON property keys).
const (
	// Content is the main text body of an entry.
	Content = "content"
	// Priority is the task priority level (high, medium, low).
	Priority = "priority"
	// Query is the search text or task number for completion.
	Query = "query"
	// Archive controls whether completed tasks are written to archive.
	Archive = "archive"
	// RecentAction is a description of what was just done (task nudge).
	RecentAction = "recent_action"
	// Caller identifies the MCP client (cursor, vscode, etc.).
	Caller = "caller"
	// Limit is the maximum number of results to return.
	Limit = "limit"
	// Since is an ISO date filter for session recall.
	Since = "since"
	// AttrFile is the metadata key on PendingUpdate recording which
	// context file was written to (e.g., "DECISIONS.md").
	AttrFile = "file"
)
