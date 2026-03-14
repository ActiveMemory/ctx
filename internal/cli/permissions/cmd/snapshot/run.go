//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/spf13/cobra"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
			return ctxerr.SettingsNotFound()
		}
		return ctxerr.FileRead(claude.Settings, readErr)
	}

	updated := false
	if _, statErr := os.Stat(claude.SettingsGolden); statErr == nil {
		updated = true
	}

	if writeErr := os.WriteFile(
		claude.SettingsGolden, content, fs.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(claude.SettingsGolden, writeErr)
	}

	write.SnapshotDone(cmd, updated, claude.SettingsGolden)
	return nil
}
