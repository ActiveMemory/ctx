//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"github.com/spf13/cobra"
)

// Run dispatches to list or show mode based on flags.
//
// Show mode is triggered by --show <id>, --latest, or a positional argument.
// Otherwise, list mode is used.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: positional arguments (session ID triggers show mode)
//   - opts: combined flags for both modes
//
// Returns:
//   - error: non-nil if the delegated command fails
func Run(cmd *cobra.Command, args []string, opts Opts) error {
	if opts.ShowID != "" || opts.Latest || len(args) > 0 {
		return runShow(cmd, args, opts)
	}

	return runList(cmd, opts)
}
