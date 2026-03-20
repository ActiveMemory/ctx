//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/format"
	"github.com/spf13/cobra"
)

// ResultLine prints a single backup result with optional SMB destination.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - scope: backup scope label (e.g. "project", "global").
//   - archive: archive file path.
//   - size: archive size in bytes.
//   - smbDest: optional SMB destination (empty string skips).
func ResultLine(cmd *cobra.Command, scope, archive string, size int64, smbDest string) {
	if cmd == nil {
		return
	}
	line := fmt.Sprintf(desc.TextDesc(text.DescKeyWriteBackupResult), scope, archive, format.Bytes(size))
	if smbDest != "" {
		line += fmt.Sprintf(desc.TextDesc(text.DescKeyWriteBackupSMBDest), smbDest)
	}
	cmd.Println(line)
}
