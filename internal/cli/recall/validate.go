//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import "github.com/ActiveMemory/ctx/internal/recall/parser"

// emptyMessage reports whether a message has no meaningful content
// (no text, tool uses, or tool results).
//
// Parameters:
//   - msg: Message to check
//
// Returns:
//   - bool: True if the message is empty
func emptyMessage(msg parser.Message) bool {
	return msg.Text == "" && len(msg.ToolUses) == 0 && len(msg.ToolResults) == 0
}
