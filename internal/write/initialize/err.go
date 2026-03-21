//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// ErrCtxNotInPath prints a multi-line diagnostic to stderr explaining
// that ctx is not in PATH, with installation instructions.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the output. Nil is a no-op.
func ErrCtxNotInPath(cmd *cobra.Command) {
	if cmd == nil {
		return
	}

	cmd.PrintErrln(desc.Text(text.DescKeyErrInitCtxNotInPath))
}
