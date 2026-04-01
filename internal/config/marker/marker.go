//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package marker

// HTML comment markers for parsing and generation.
const (
	// CommentOpen is the HTML comment opening tag.
	CommentOpen = "<!--"
	// CommentClose is the HTML comment closing tag.
	CommentClose = "-->"
)

// HTML element markers for preformatted content.
const (
	// TagPre is the opening tag for preformatted text blocks.
	TagPre = "<pre>"
	// TagPreClose is the closing tag for preformatted text blocks.
	TagPreClose = "</pre>"
	// TagPreCode is the compound opening tag for code within pre blocks.
	TagPreCode = "<pre><code>"
	// TagPreCodeClose is the compound closing tag.
	TagPreCodeClose = "</code></pre>"
	// TagDetails is the opening tag for collapsible details blocks.
	TagDetails = "<details>"
	// TagDetailsClose is the closing tag for details blocks.
	TagDetailsClose = "</details>"
	// TagSummaryOpen is the opening tag for summary elements.
	TagSummaryOpen = "<summary>"
	// TagSummaryClose is the closing tag for summary elements.
	TagSummaryClose = "</summary>"
)

// HTML entity escapes.
const (
	// EntityLT is the HTML entity for <.
	EntityLT = "&lt;"
	// EntityGT is the HTML entity for >.
	EntityGT = "&gt;"
	// AngleLT is the literal < character.
	AngleLT = "<"
	// AngleGT is the literal > character.
	AngleGT = ">"
)

// Markdown inline formatting.
const (
	// BoldWrap is the Markdown bold delimiter.
	BoldWrap = "**"
)

// Context block markers for embedding context in files.
const (
	// CtxStart marks the beginning of an embedded context block.
	CtxStart = "<!-- ctx:context -->"
	// CtxEnd marks the end of an embedded context block.
	CtxEnd = "<!-- ctx:end -->"
)

// Copilot block markers for .github/copilot-instructions.md.
const (
	// CopilotStart marks the beginning of ctx-managed Copilot content.
	CopilotStart = "<!-- ctx:copilot -->"
	// CopilotEnd marks the end of ctx-managed Copilot content.
	CopilotEnd = "<!-- ctx:copilot:end -->"
)

// Agents block markers for AGENTS.md.
const (
	// AgentsStart marks the beginning of ctx-managed AGENTS.md content.
	AgentsStart = "<!-- ctx:agents -->"
	// AgentsEnd marks the end of ctx-managed AGENTS.md content.
	AgentsEnd = "<!-- ctx:agents:end -->"
)

// Index markers for auto-generated table of contents sections.
const (
	// IndexStart marks the beginning of an auto-generated index.
	IndexStart = "<!-- INDEX:START -->"
	// IndexEnd marks the end of an auto-generated index.
	IndexEnd = "<!-- INDEX:END -->"
)

// Task checkbox prefixes for Markdown task lists.
const (
	// PrefixTaskUndone is the prefix for an unchecked task item.
	PrefixTaskUndone = "- [ ]"
	// PrefixTaskDone is the prefix for a checked (completed) task item.
	PrefixTaskDone = "- [x]"
)

const (
	// MarkTaskComplete is the unchecked task marker.
	MarkTaskComplete = "x"
)

// Published block markers for MEMORY.md.
const (
	// PublishStart begins the ctx-published block in MEMORY.md.
	PublishStart = "<!-- ctx:published -->"
	// PublishEnd ends the ctx-published block in MEMORY.md.
	PublishEnd = "<!-- ctx:end -->"
)

// Entry status markers for knowledge files.
const (
	// PrefixSuperseded is the strikethrough prefix that marks an entry as
	// superseded by a newer entry.
	PrefixSuperseded = "~~Superseded"
)

// System reminder tags injected by Claude Code into tool results.
const (
	// TagSystemReminderOpen is the opening tag for system reminders.
	TagSystemReminderOpen = "<system-reminder>"
	// TagSystemReminderClose is the closing tag for system reminders.
	TagSystemReminderClose = "</system-reminder>"
)

// Context compaction artifacts injected by Claude Code when the conversation
// exceeds the context window. The compaction injects two blocks as a user
// message in the JSONL transcript:
//
//  1. A multi-line <summary>...</summary> block containing a structured
//     conversation summary (sections: Request/Intent, Technical Concepts,
//     Files, Current State).
//  2. A boilerplate continuation prompt ("If you need specific details
//     from before compaction...read the full transcript at...").
//
// INVARIANT: Claude Code's <summary> tag always appears alone on its own
// line (the content is inherently multi-line). Our <summary> tags are always
// single-line (<summary>N lines</summary>, see TplRecallDetailsOpen).
// This invariant makes disambiguation safe: a line that is exactly "<summary>"
// is a compaction artifact; a line containing "<summary>...</summary>" is ours.
const (
	// TagCompactionSummaryOpen is a standalone <summary> on its own line.
	TagCompactionSummaryOpen = "<summary>"
	// TagCompactionSummaryClose is the closing </summary> tag.
	TagCompactionSummaryClose = "</summary>"
	// CompactionBoilerplatePrefix starts the continuation prompt after
	// a compaction summary.
	CompactionBoilerplatePrefix = "If you need specific details from before compaction"
)

// Markdown table markers for index generation.
const (
	// TablePipe is the cell delimiter in Markdown tables.
	TablePipe = "|"
	// TablePipePad is the padded cell delimiter.
	TablePipePad = " | "
	// TableRowOpen opens a table row.
	TableRowOpen = "| "
	// TableRowClose closes a table row.
	TableRowClose = " |"
	// TableSepCell is a header separator cell.
	TableSepCell = "------"
	// TablePipeEscaped is an escaped pipe for use inside cells.
	TablePipeEscaped = "\\|"
	// TableRowFmt is the Printf format for a two-column table row.
	TableRowFmt = "| %s | %s |"
	// TableSepFmt is the Printf format for a two-column separator row.
	TableSepFmt = "|%s|%s|"
)
