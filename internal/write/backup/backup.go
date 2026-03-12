//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/write"
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// BackupResultLine prints a single backup result with optional SMB destination.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - scope: backup scope label (e.g. "project", "global").
//   - archive: archive file path.
//   - size: archive size in bytes.
//   - smbDest: optional SMB destination (empty string skips).
func BackupResultLine(cmd *cobra.Command, scope, archive string, size int64, smbDest string) {
	if cmd == nil {
		return
	}
	line := fmt.Sprintf(config.TplBackupResult, scope, archive, write.FormatBytes(size))
	if smbDest != "" {
		line += fmt.Sprintf(config.TplBackupSMBDest, smbDest)
	}
	cmd.Println(line)
}
