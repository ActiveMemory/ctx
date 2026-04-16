//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize resolves and canonicalizes section
// names for the add command.
//
// # TargetSection
//
// [TargetSection] ensures a user-provided section name
// carries the correct Markdown heading prefix. If the
// input does not already start with "## ", the function
// prepends it. The heading prefix is read from the
// token.HeadingLevelThreeStart constant.
//
// Callers must not pass an empty string. The empty case
// is handled by insert.Task before this function is
// reached, so TargetSection always receives a non-empty
// section name.
//
// # Usage Example
//
// A user runs:
//
//	ctx add task "fix tests" --section "Phase 1"
//
// The insert subpackage calls TargetSection("Phase 1")
// which returns "## Phase 1". The insert logic then
// searches for that heading in TASKS.md and places the
// new task below it.
//
// # Data Flow
//
// The insert subpackage is the sole caller. It invokes
// TargetSection inside TaskAfterSection before scanning
// the file content for the heading. This keeps heading
// format knowledge centralized in one function.
package normalize
