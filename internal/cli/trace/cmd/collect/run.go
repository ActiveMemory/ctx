//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// Run collects context refs from all sources and outputs the trailer to stdout.
//
// Parameters:
//   - cmd: Cobra command for output stream
//
// Returns:
//   - error: non-nil on execution failure
func Run(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	refs := trace.Collect(contextDir)
	trailer := trace.FormatTrailer(refs)
	if trailer != "" {
		cmd.Println(trailer)
	}
	return nil
}
