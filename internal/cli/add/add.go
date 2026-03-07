//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/cmd/root"
	"github.com/ActiveMemory/ctx/internal/cli/add/core"
)

// EntryParams is the public type for entry parameters, re-exported from core.
type EntryParams = core.EntryParams

// ValidateEntry checks required fields for a given entry type.
var ValidateEntry = root.ValidateEntry

// WriteEntry formats and writes an entry to the appropriate context file.
var WriteEntry = root.WriteEntry

// Cmd returns the "ctx add" command.
//
// Returns:
//   - *cobra.Command: Configured add command
func Cmd() *cobra.Command {
	return root.Cmd()
}
