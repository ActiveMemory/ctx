//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package change

import (
	"github.com/spf13/cobra"
)

// List prints rendered change output.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Pre-rendered changes string.
func List(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	cmd.Print(content)
}
