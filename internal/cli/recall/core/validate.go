//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/err/journal"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// EmptyMessage reports whether a message has no meaningful content
// (no text, tool uses, or tool results).
//
// Parameters:
//   - msg: Message to check
//
// Returns:
//   - bool: True if the message is empty
func EmptyMessage(msg parser.Message) bool {
	return msg.Text == "" && len(msg.ToolUses) == 0 && len(msg.ToolResults) == 0
}

// ValidateExportFlags checks for invalid flag combinations.
//
// Parameters:
//   - args: positional arguments (session IDs).
//   - opts: export flag values.
//
// Returns:
//   - error: non-nil if flags conflict.
func ValidateExportFlags(args []string, opts ExportOpts) error {
	if len(args) > 0 && opts.All {
		return ctxErr.AllWithID()
	}
	if opts.Regenerate && !opts.All {
		return journal.RegenerateRequiresAll()
	}
	return nil
}
