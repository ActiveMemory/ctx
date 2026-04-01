//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

// Claude Code tool input JSON keys for display formatting.
const (
	ToolInputFilePath    = "file_path"
	ToolInputCommand     = "command"
	ToolInputPattern     = "pattern"
	ToolInputURL         = "url"
	ToolInputQuery       = "query"
	ToolInputDescription = "description"
)

// Tool display limits.
const (
	// ToolDisplayMaxLen is the max length for tool parameter
	// display before truncation.
	ToolDisplayMaxLen = 100
)
