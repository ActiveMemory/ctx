//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/spf13/cobra"
)

// ShowContent prints prompt template content to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Raw prompt template bytes.
func ShowContent(cmd *cobra.Command, content []byte) {
	if cmd == nil {
		return
	}
	cmd.Print(string(content))
}
