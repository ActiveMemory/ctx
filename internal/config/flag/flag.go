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

// CLI flag prefixes for display formatting.
const PrefixLong = "--"

// Add command flag names: used for both flag registration and error display.
const (
	Context      = "context"
	Rationale    = "rationale"
	Consequences = "consequences"
	Lesson       = "lesson"
	Application  = "application"
)
