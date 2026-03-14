//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// taskPattern captures indent, checkbox state, and content.
//
// Pattern: ^(\s*)-\s*\[([x ]?)]\s*(.+)$
//
// Groups:
//   - 1: indent (leading whitespace, may be empty)
//   - 2: state ("x" for completed, " " or "" for pending)
//   - 3: content (task text)
const taskPattern = `^(\s*)-\s*\[([x ]?)]\s*(.+)$`

// Task matches a task item on a single line.
//
// Use with MatchString or FindStringSubmatch on individual lines.
// For multiline content, use TaskMultiline.
var Task = regexp.MustCompile(taskPattern)

// TaskMultiline matches task items across multiple lines.
//
// Use with FindAllStringSubmatch on multiline content.
var TaskMultiline = regexp.MustCompile(`(?m)` + taskPattern)

// TaskDoneTimestamp extracts the #done: timestamp from a task line.
//
// Groups:
//   - 1: timestamp (YYYY-MM-DD-HHMMSS)
var TaskDoneTimestamp = regexp.MustCompile(`#done:(\d{4}-\d{2}-\d{2}-\d{6})`)

// Runtime configuration.
const (
	// TaskCompleteReplace is the regex replacement string for marking a task done.
	TaskCompleteReplace = "$1- [x] $3"
)
