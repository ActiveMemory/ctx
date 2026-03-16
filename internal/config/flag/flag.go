//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flag

// Global CLI flag names.
const (
	ContextDir      = "context-dir"
	AllowOutsideCwd = "allow-outside-cwd"
)

// PrefixLong is a CLI flag prefix for display formatting.
const PrefixLong = "--"

// Add command flag names: used for both flag registration and error display.
const (
	Context     = "context"
	Rationale   = "rationale"
	Consequence = "consequence"
	Lesson      = "lesson"
	Application = "application"
)
