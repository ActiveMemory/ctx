//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sprintf formats a string and prints it to the command's stdout stream.
//
// This is the internal building block for all formatted output in the package.
//
// Parameters:
//   - cmd: Cobra command whose stdout stream receives the output.
//   - format: fmt.Sprintf format string.
//   - args: format arguments.
func sprintf(cmd *cobra.Command, format string, args ...any) {
	cmd.Println(fmt.Sprintf(format, args...))
}
