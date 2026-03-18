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

	claude2 "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/err/config"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	errparser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/write/restore"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/cli/permissions/core"
)

// Run resets settings.local.json from the golden image.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write/parse failure or missing golden file
func Run(cmd *cobra.Command) error {
	goldenBytes, goldenReadErr := os.ReadFile(claude2.SettingsGolden)
	if goldenReadErr != nil {
		if os.IsNotExist(goldenReadErr) {
			return config.GoldenNotFound()
		}
		return fs2.FileRead(claude2.SettingsGolden, goldenReadErr)
	}

	localBytes, localReadErr := os.ReadFile(claude2.Settings)
	if localReadErr != nil {
		if os.IsNotExist(localReadErr) {
			if writeErr := os.WriteFile(
				claude2.Settings, goldenBytes, fs.PermFile,
			); writeErr != nil {
				return fs2.FileWrite(claude2.Settings, writeErr)
			}
			restore.RestoreNoLocal(cmd)
			return nil
		}
		return fs2.FileRead(claude2.Settings, localReadErr)
	}

	if bytes.Equal(goldenBytes, localBytes) {
		restore.RestoreMatch(cmd)
		return nil
	}

	var golden, local claude.Settings
	if goldenParseErr := json.Unmarshal(goldenBytes, &golden); goldenParseErr != nil {
		return errparser.ParseFile(claude2.SettingsGolden, goldenParseErr)
	}
	if localParseErr := json.Unmarshal(localBytes, &local); localParseErr != nil {
		return errparser.ParseFile(claude2.Settings, localParseErr)
	}

	restored, dropped := core.DiffStringSlices(
		golden.Permissions.Allow, local.Permissions.Allow,
	)
	denyRestored, denyDropped := core.DiffStringSlices(
		golden.Permissions.Deny, local.Permissions.Deny,
	)

	restore.RestoreDiff(cmd, dropped, restored, denyDropped, denyRestored)

	if writeErr := os.WriteFile(
		claude2.Settings, goldenBytes, fs.PermFile,
	); writeErr != nil {
		return fs2.FileWrite(claude2.Settings, writeErr)
	}

	restore.RestoreDone(cmd)
	return nil
}
