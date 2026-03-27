//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/err/journal"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
)

// EmptyMessage reports whether a message has no meaningful content
// (no text, tool uses, or tool results).
//
// Parameters:
//   - msg: Message to check
//
// Returns:
//   - bool: True if the message is empty
func EmptyMessage(msg entity.Message) bool {
	return msg.Text == "" && len(msg.ToolUses) == 0 && len(msg.ToolResults) == 0
}

// ImportFlags checks for invalid flag combinations.
//
// Parameters:
//   - args: positional arguments (session IDs).
//   - opts: import flag values.
//
// Returns:
//   - error: non-nil if flags conflict.
func ImportFlags(args []string, opts entity.ImportOpts) error {
	if len(args) > 0 && opts.All {
		return errSession.AllWithID()
	}
	if opts.Regenerate && !opts.All {
		return journal.RegenerateRequiresAll()
	}
	return nil
}
