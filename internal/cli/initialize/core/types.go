//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "github.com/spf13/cobra"

// MergeParams holds configuration for the create-or-merge file operation.
type MergeParams struct {
	Filename        string
	MarkerStart     string
	TemplateContent []byte
	Force           bool
	AutoMerge       bool
	ConfirmPrompt   string
	UpdateFn        func(cmd *cobra.Command, existing string, tpl []byte) error
}
