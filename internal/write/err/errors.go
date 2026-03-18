//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// WithError writes a prefixed error message to the command's stderr stream.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the message. Nil is a no-op.
//   - err: the error to display after the "Error: " prefix.
func WithError(cmd *cobra.Command, err error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(assets.TextDesc(assets.TextDescKeyWritePrefixError), err)
}

// WarnFile prints a non-fatal file operation warning to stderr.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the message. Nil is a no-op.
//   - path: path of the file that caused the warning.
//   - err: the underlying error.
func WarnFile(cmd *cobra.Command, path string, err error) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf("  ! %s: %v", path, err))
}
