//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/permissions/core"
	"github.com/ActiveMemory/ctx/internal/config"
)

// Run saves settings.local.json as the golden image.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write failure or missing settings file
func Run(cmd *cobra.Command) error {
	content, readErr := os.ReadFile(config.FileSettings)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return core.ErrSettingsNotFound()
		}
		return core.ErrReadFile(config.FileSettings, readErr)
	}

	// Determine message based on whether golden already exists.
	verb := "Saved"
	if _, statErr := os.Stat(config.FileSettingsGolden); statErr == nil {
		verb = "Updated"
	}

	if writeErr := os.WriteFile(config.FileSettingsGolden, content, config.PermFile); writeErr != nil {
		return core.ErrWriteFile(config.FileSettingsGolden, writeErr)
	}

	cmd.Println(fmt.Sprintf("%s golden image: %s", verb, config.FileSettingsGolden))
	return nil
}
