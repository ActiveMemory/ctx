//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// FileTokenEntry tracks per-file token counts during context injection.
//
// Fields:
//   - Name: Context file name
//   - Tokens: Estimated token count
type FileTokenEntry struct {
	Name   string
	Tokens int
}

// MessageListEntry holds the data for a single row in the message list output.
//
// Fields:
//   - Hook: Hook lifecycle event name
//   - Variant: Message variant within the hook
//   - Category: Message category (nudge, relay, block, etc.)
//   - Description: Human-readable description
//   - TemplateVars: Variable names used in the template
//   - HasOverride: Whether a user override exists
type MessageListEntry struct {
	Hook         string   `json:"hook"`
	Variant      string   `json:"variant"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	TemplateVars []string `json:"template_vars"`
	HasOverride  bool     `json:"has_override"`
}

// StaleEntry describes a file that has not been modified within the
// configured freshness window.
type StaleEntry struct {
	// Path is the relative file path.
	Path string
	// Desc is the human-readable file description.
	Desc string
	// ReviewURL is the optional URL for reviewing the file against upstream.
	ReviewURL string
	// Days is the number of days since last modification.
	Days int
}
