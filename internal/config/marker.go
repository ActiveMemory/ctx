//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// HTML comment markers for parsing and generation.
const (
	// CommentOpen is the HTML comment opening tag.
	CommentOpen = "<!--"
	// CommentClose is the HTML comment closing tag.
	CommentClose = "-->"
)

// Context block markers for embedding context in files.
const (
	// CtxMarkerStart marks the beginning of an embedded context block.
	CtxMarkerStart = "<!-- ctx:context -->"
	// CtxMarkerEnd marks the end of an embedded context block.
	CtxMarkerEnd = "<!-- ctx:end -->"
)

// Prompt block markers for PROMPT.md.
const (
	// PromptMarkerStart marks the beginning of the prompt block.
	PromptMarkerStart = "<!-- ctx:prompt -->"
	// PromptMarkerEnd marks the end of the prompt block.
	PromptMarkerEnd = "<!-- ctx:prompt:end -->"
)

// Plan block markers for IMPLEMENTATION_PLAN.md.
const (
	// PlanMarkerStart marks the beginning of the plan block.
	PlanMarkerStart = "<!-- ctx:plan -->"
	// PlanMarkerEnd marks the end of the plan block.
	PlanMarkerEnd = "<!-- ctx:plan:end -->"
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
