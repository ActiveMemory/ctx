//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// MergeParams holds configuration for create-or-merge file operations.
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
