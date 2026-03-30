//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// MergeParams holds configuration for create-or-merge file operations.
//
// Fields:
//   - Filename: Target file name
//   - MarkerStart: Merge region start marker
//   - MarkerEnd: Merge region end marker
//   - TemplateContent: Embedded template bytes to merge
//   - Force: Overwrite without merge
//   - AutoMerge: Merge without user confirmation
//   - ConfirmPrompt: Prompt text for interactive confirmation
//   - UpdateTextKey: Text key for the update message
type MergeParams struct {
	Filename        string
	MarkerStart     string
	MarkerEnd       string
	TemplateContent []byte
	Force           bool
	AutoMerge       bool
	ConfirmPrompt   string
	UpdateTextKey   string
}
