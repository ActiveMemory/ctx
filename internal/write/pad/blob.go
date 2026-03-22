//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// BlobWritten prints confirmation that a blob was written to a file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - size: number of bytes written.
//   - path: output file path.
func BlobWritten(cmd *cobra.Command, size int, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadBlobWritten), size, path),
	)
}

// BlobShow prints raw blob data to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - data: Raw blob bytes.
func BlobShow(cmd *cobra.Command, data []byte) {
	if cmd == nil {
		return
	}
	cmd.Print(string(data))
}
