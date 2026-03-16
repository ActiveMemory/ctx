//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// EntryField pairs a label TextDescKey with a Value from prompt
// arguments.
type EntryField struct {
	LabelKey string
	Value    string
}

// EntryPromptSpec holds the TextDescKeys that vary between
// add-decision and add-learning prompts.
type EntryPromptSpec struct {
	HeaderKey  string
	FooterKey  string
	FieldFmtK  string
	ResultDKey string
	Fields     []EntryField
}
