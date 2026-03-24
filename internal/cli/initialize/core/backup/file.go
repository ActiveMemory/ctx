//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"
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
	if writeErr := os.WriteFile(backupName, content, fs.PermFile); writeErr != nil {
		return errBackup.Create(backupName, writeErr)
	}
	initialize.Backup(cmd, backupName)
	return nil
}
