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
