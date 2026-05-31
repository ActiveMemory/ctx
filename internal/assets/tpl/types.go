//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// ObsidianData is the render data for [ObsidianReadme].
type ObsidianData struct {
	// JournalDir is the journal source directory path.
	JournalDir string
}

// JournalSiteData is the render data for [JournalSiteReadme].
type JournalSiteData struct {
	// JournalDir is the journal source directory path.
	JournalDir string
}

// TriggerData is the render data for [TriggerScript].
type TriggerData struct {
	// Name is the trigger script base name (without .sh).
	Name string
	// Type is the trigger type (e.g. pre-tool-use, session-start).
	Type string
}

// LearningData is the render data for [Learning].
type LearningData struct {
	// Timestamp is the entry creation timestamp.
	Timestamp string
	// Title is the learning title/summary.
	Title string
	// Context is what prompted the learning.
	Context string
	// Lesson is the key insight.
	Lesson string
	// Application is how to apply it going forward.
	Application string
}

// LoopData is the render data for [LoopScript].
type LoopData struct {
	// PromptFile is the absolute path to the loop's prompt file.
	PromptFile string
	// CompletionSignal is the string that, when seen in tool output,
	// ends the loop.
	CompletionSignal string
	// MaxIter is the iteration cap; 0 means unlimited (the
	// iteration-limit block is omitted).
	MaxIter int
	// AICommand is the shell command that runs the AI tool.
	AICommand string
	// LoopComplete is the completion banner line.
	LoopComplete string
}

// DecisionData is the render data for [Decision].
type DecisionData struct {
	// Timestamp is the entry creation timestamp.
	Timestamp string
	// Title is the decision title/summary.
	Title string
	// Context is what prompted the decision.
	Context string
	// Rationale is why this choice over alternatives.
	Rationale string
	// Consequence is what changes as a result.
	Consequence string
}

// MetaTableData is the render data for [MetaTable].
type MetaTableData struct {
	// Summary is the <summary> text for the collapsible block.
	Summary string
	// Rows are the table's label/value rows, in order.
	Rows []MetaRow
}

// MetaRow is one label/value row in a [MetaTable].
type MetaRow struct {
	// Label is the row's bold left-column text.
	Label string
	// Value is the row's right-column text.
	Value string
}

// DetailsData is the render data for [Details].
type DetailsData struct {
	// Summary is the <summary> text for the collapsible block.
	Summary string
	// Body is the pre-rendered block body (already escaped/wrapped by
	// the caller).
	Body string
}
