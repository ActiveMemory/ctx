//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
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

	cmd.PrintErrln("Error: ctx is not in your PATH")
	cmd.PrintErrln(
		"The hooks created by 'ctx init' require ctx to be in your PATH.",
	)
	cmd.PrintErrln("Without this, Claude Code hooks will fail silently.")
	cmd.PrintErrln()
	cmd.PrintErrln("To fix this:")
	cmd.PrintErrln("  1. Build:   make build")
	cmd.PrintErrln("  2. Install: sudo make install")
	cmd.PrintErrln()
	cmd.PrintErrln("Or manually:")
	cmd.PrintErrln("  sudo cp ./ctx /usr/local/bin/")
	cmd.PrintErrln()
	cmd.PrintErrln("Then run 'ctx init' again.")
}
