//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"github.com/spf13/cobra"

	journalImporter "github.com/ActiveMemory/ctx/internal/cli/journal/cmd/importer"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// Run delegates to the journal importer. The journal package is canonical;
// this wrapper exists for backward compatibility with recall command wiring.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: positional arguments (optional session ID)
//   - opts: import flag values
//
// Returns:
//   - error: non-nil on validation, scan, or write failures
func Run(cmd *cobra.Command, args []string, opts entity.ImportOpts) error {
	return journalImporter.Run(cmd, args, opts)
}
