//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"github.com/spf13/cobra"

	coreEdit "github.com/ActiveMemory/ctx/internal/cli/pad/core/edit"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run edits a scratchpad entry based on the selected mode.
//
// Parameters:
//   - cmd: Cobra command for output
//   - opts: Edit operation parameters
//
// Returns:
//   - error: Non-nil on invalid index, type mismatch, or read/write failure
func Run(cmd *cobra.Command, opts Opts) error {
	var entries []string
	var editErr error

	switch opts.Mode {
	case ModeAppend:
		entries, editErr = coreEdit.Append(opts.N, opts.Text)
	case ModePrepend:
		entries, editErr = coreEdit.Prepend(opts.N, opts.Text)
	case ModeBlob:
		entries, editErr = coreEdit.UpdateBlob(
			opts.N, opts.FilePath, opts.LabelText,
		)
	default:
		entries, editErr = coreEdit.Replace(opts.N, opts.Text)
	}
	if editErr != nil {
		return editErr
	}

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.EntryUpdated(cmd, opts.N)
	return nil
}
