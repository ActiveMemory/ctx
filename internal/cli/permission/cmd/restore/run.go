//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/cli/permission/core"
	configClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/write/restore"
)

// Run resets settings.local.json from the golden image.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write/parse failure or missing golden file
func Run(cmd *cobra.Command) error {
	goldenBytes, goldenReadErr := os.ReadFile(configClaude.SettingsGolden)
	if goldenReadErr != nil {
		if os.IsNotExist(goldenReadErr) {
			return config.GoldenNotFound()
		}
		return errFs.FileRead(configClaude.SettingsGolden, goldenReadErr)
	}

	localBytes, localReadErr := os.ReadFile(configClaude.Settings)
	if localReadErr != nil {
		if os.IsNotExist(localReadErr) {
			if writeErr := os.WriteFile(
				configClaude.Settings, goldenBytes, fs.PermFile,
			); writeErr != nil {
				return errFs.FileWrite(configClaude.Settings, writeErr)
			}
			restore.RestoreNoLocal(cmd)
			return nil
		}
		return errFs.FileRead(configClaude.Settings, localReadErr)
	}

	if bytes.Equal(goldenBytes, localBytes) {
		restore.RestoreMatch(cmd)
		return nil
	}

	var golden, local claude.Settings
	if goldenParseErr := json.Unmarshal(goldenBytes, &golden); goldenParseErr != nil {
		return errParser.ParseFile(configClaude.SettingsGolden, goldenParseErr)
	}
	if localParseErr := json.Unmarshal(localBytes, &local); localParseErr != nil {
		return errParser.ParseFile(configClaude.Settings, localParseErr)
	}

	restored, dropped := core.DiffStringSlices(
		golden.Permissions.Allow, local.Permissions.Allow,
	)
	denyRestored, denyDropped := core.DiffStringSlices(
		golden.Permissions.Deny, local.Permissions.Deny,
	)

	restore.RestoreDiff(cmd, dropped, restored, denyDropped, denyRestored)

	if writeErr := os.WriteFile(
		configClaude.Settings, goldenBytes, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(configClaude.Settings, writeErr)
	}

	restore.RestoreDone(cmd)
	return nil
}
