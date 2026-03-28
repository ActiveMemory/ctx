//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// CheckResult holds the outcome of checking a stage value.
type CheckResult struct {
	Value string
}

// MarkResult holds the outcome of marking a stage.
type MarkResult struct {
	Marked bool
}
