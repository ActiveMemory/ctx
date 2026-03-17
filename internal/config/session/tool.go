//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

// Tool identifiers for session parsers.
const (
	// ToolClaudeCode is the tool identifier for Claude Code sessions.
	ToolClaudeCode = "claude-code"
	// ToolCopilot is the tool identifier for VS Code Copilot Chat sessions.
	ToolCopilot = "copilot"
	// ToolMarkdown is the tool identifier for Markdown session files.
	ToolMarkdown = "markdown"
)

// Claude Code tool names used in session transcripts.
const (
	ToolRead      = "Read"
	ToolWrite     = "Write"
	ToolEdit      = "Edit"
	ToolBash      = "Bash"
	ToolGrep      = "Grep"
	ToolGlob      = "Glob"
	ToolWebFetch  = "WebFetch"
	ToolWebSearch = "WebSearch"
	ToolTask      = "Task"
)
