//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"
)

// InfoExistsWritingAsAlternative reports that a file already exists and the
// content is being written to an alternative filename instead.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: the original target path that already exists.
//   - alternative: the fallback path where content was written.
func InfoExistsWritingAsAlternative(
	cmd *cobra.Command, path, alternative string,
) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteExistsWritingAsAlternative), path, alternative))
}
