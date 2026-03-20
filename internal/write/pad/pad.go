//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// InfoPathConversionExists reports that a pad export target already
// exists and will be written with a timestamped alternative name.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - rootDir: export destination directory.
//   - oldPath: original blob label.
//   - newPath: timestamped alternative name joined with rootDir.
func InfoPathConversionExists(
	cmd *cobra.Command, rootDir, oldPath, newPath string,
) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.TextDesc(text.DescKeyWritePathExists), oldPath, filepath.Join(rootDir, newPath),
		),
	)
}
