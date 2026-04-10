//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// PromptEntryField pairs a label DescKey with a Value from prompt
// arguments.
//
// Fields:
//   - KeyLabel: Text DescKey for the field label
//   - Value: User-provided argument value
type PromptEntryField struct {
	KeyLabel string
	Value    string
}

// PromptEntrySpec holds the DescKeys that vary between
// add-decision and add-learning MCP prompts.
//
// Fields:
//   - KeyHeader: DescKey for the prompt header
//   - KeyFooter: DescKey for the prompt footer
//   - FieldFmtK: DescKey format string for field labels
//   - KeyResultD: DescKey for the result description
//   - Fields: Variable fields for this entry type
type PromptEntrySpec struct {
	KeyHeader  string
	KeyFooter  string
	FieldFmtK  string
	KeyResultD string
	Fields     []PromptEntryField
}
