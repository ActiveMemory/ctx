//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/spf13/cobra"
)

// ErrWithError writes a prefixed error message to the command's stderr stream.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the message. Nil is a no-op.
//   - err: the error to display after the "Error: " prefix.
func ErrWithError(cmd *cobra.Command, err error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(prefixError, err)
}

// WarnFileErr prints a non-fatal file operation warning to stderr.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the message. Nil is a no-op.
//   - path: path of the file that caused the warning.
//   - err: the underlying error.
func WarnFileErr(cmd *cobra.Command, path string, err error) {
	if cmd == nil {
		return
	}
	sprintf(cmd, "  ! %s: %v", path, err)
}
