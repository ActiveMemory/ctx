//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/err/config"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/restore"
)

// Run saves settings.local.json as the golden image.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write failure or missing settings file
func Run(cmd *cobra.Command) error {
	content, readErr := os.ReadFile(claude.Settings)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return config.SettingsNotFound()
		}
		return ctxErr.FileRead(claude.Settings, readErr)
	}

	updated := false
	if _, statErr := os.Stat(claude.SettingsGolden); statErr == nil {
		updated = true
	}

	if writeErr := os.WriteFile(
		claude.SettingsGolden, content, fs.PermFile,
	); writeErr != nil {
		return ctxErr.FileWrite(claude.SettingsGolden, writeErr)
	}

	restore.SnapshotDone(cmd, updated, claude.SettingsGolden)
	return nil
}
