//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// File creates a timestamped .bak copy and reports it.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Original file path
//   - content: Content to back up
//
// Returns:
//   - error: Non-nil if the backup write fails
func File(cmd *cobra.Command, filename string, content []byte) error {
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf(file.BackupFormat, filename, timestamp)
	if writeErr := ctxIo.SafeWriteFile(
		backupName, content, fs.PermFile,
	); writeErr != nil {
		return errBackup.Create(backupName, writeErr)
	}
	initialize.Backup(cmd, backupName)
	return nil
}
