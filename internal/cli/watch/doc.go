//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package watch implements the "ctx watch" command for processing
// structured context-update commands from AI output.
//
// The watch command reads stdin for XML-style update commands
// (task, decision, learning, convention, complete) and applies them
// to the corresponding context files. It supports dry-run mode
// for previewing changes.
package watch
