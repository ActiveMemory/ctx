//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// Bold metadata field prefixes in journal/session Markdown.
const (
	// MetadataID is the bold ID field prefix.
	MetadataID = "**ID**:"
	// MetadataDate is the bold date field prefix.
	MetadataDate = "**Date**:"
	// MetadataTime is the bold time field prefix.
	MetadataTime = "**Time**:"
	// MetadataDuration is the bold duration field prefix.
	MetadataDuration = "**Duration**:"
	// MetadataTool is the bold tool field prefix.
	MetadataTool = "**Tool**:"
	// MetadataProject is the bold project field prefix.
	MetadataProject = "**Project**:"
	// MetadataBranch is the bold branch field prefix.
	MetadataBranch = "**Branch**:"
	// MetadataModel is the bold model field prefix.
	MetadataModel = "**Model**:"
	// MetadataTurns is the bold turns field prefix.
	MetadataTurns = "**Turns**:"
	// MetadataParts is the bold parts field prefix.
	MetadataParts = "**Parts**:"
	// MetadataType is the bold type field prefix.
	MetadataType = "**Type**:"
	// MetadataStartTime is the bold start_time field prefix.
	MetadataStartTime = "**start_time**:"
	// MetadataEndTime is the bold end_time field prefix.
	MetadataEndTime = "**end_time**:"
	// MetadataSource is the bold source field prefix.
	MetadataSource = "**Source**:"
)

// Conversation role display labels used in exported journal entries.
const (
	// LabelRoleUser is the display label for user turns.
	LabelRoleUser = "User"
	// LabelRoleAssistant is the display label for assistant turns.
	LabelRoleAssistant = "Assistant"
)

// Journal content markers for detecting session modes.
const (
	// LabelSuggestionMode identifies suggestion mode sessions in journal content.
	LabelSuggestionMode = "SUGGESTION MODE:"
)

// Journal turn markers for content transformation.
const (
	// LabelBoldReminder is the bold-style system reminder prefix.
	LabelBoldReminder = "**System Reminder**:"
	// LabelToolOutput is the turn role label for tool output turns.
	LabelToolOutput = "Tool Output"
)
